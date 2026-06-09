package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"linkworld-backend/internal/blockchain"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T6 组批结算 TDD（design §6.1 编排 + §7.0① L2 + arch-review B1-L2 冷启动回退红线）。
//
// 用 SettlementClient 接口 mock 链交互（避开 T3 的 Cancun/simulated.Backend 限制；
// T6 测组批/熔断/幂等编排逻辑，非链交互本身——L3 发交易前断言已在 client_test 覆盖）。

// fakeSettlementClient 记录每次 MonthlySettlement 调用的批次，按预设结果返回 receipt/error。
type fakeSettlementClient struct {
	mu    sync.Mutex
	calls []fakeBatchCall
	// resultFn 决定第 callIdx 次调用（0-based）返回的 receipt 与 error。nil → 默认全部成功(status=1)。
	resultFn func(callIdx int, users []common.Address) (*gethtypes.Receipt, error)
}

type fakeBatchCall struct {
	users       []common.Address
	operatorIds []*big.Int
	amounts     []*big.Int
}

func (f *fakeSettlementClient) MonthlySettlement(ctx context.Context, users []common.Address, operatorIds, amounts []*big.Int) (*gethtypes.Receipt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	idx := len(f.calls)
	f.calls = append(f.calls, fakeBatchCall{users: users, operatorIds: operatorIds, amounts: amounts})
	if f.resultFn != nil {
		return f.resultFn(idx, users)
	}
	return &gethtypes.Receipt{Status: gethtypes.ReceiptStatusSuccessful}, nil
}

// inMemBatchStore 是 SettlementBatchStore 的内存测试替身（避开 gorm/sqlite，聚焦编排逻辑）。
type inMemBatchStore struct {
	mu sync.Mutex
	// key: month|batchIndex
	batches map[string]*settlementBatchRecord
	// historyTotals 模拟历史月份各批总额（供月均熔断计算）。
	historyTotals []*big.Int
}

func newInMemBatchStore() *inMemBatchStore {
	return &inMemBatchStore{batches: map[string]*settlementBatchRecord{}}
}

func bkey(month string, idx int) string { return fmt.Sprintf("%s|%d", month, idx) }

func (s *inMemBatchStore) GetBatch(month string, batchIndex int) (*settlementBatchRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.batches[bkey(month, batchIndex)]; ok {
		cp := *r
		return &cp, nil
	}
	return nil, nil
}

func (s *inMemBatchStore) SaveBatch(r *settlementBatchRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *r
	s.batches[bkey(r.Month, r.BatchIndex)] = &cp
	return nil
}

func (s *inMemBatchStore) HistoricalBatchTotals(beforeMonth string) ([]*big.Int, error) {
	return s.historyTotals, nil
}

// makeUsers 造 n 个确定地址。
func makeUsers(n int) []common.Address {
	out := make([]common.Address, n)
	for i := 0; i < n; i++ {
		out[i] = common.BigToAddress(big.NewInt(int64(i + 1)))
	}
	return out
}

// makeItems 造 n 条结算项，每项金额 amountEach（6 位最小单位）。
func makeItems(n int, amountEach int64) []SettlementItem {
	users := makeUsers(n)
	out := make([]SettlementItem, n)
	for i := 0; i < n; i++ {
		out[i] = SettlementItem{
			User:       users[i],
			OperatorID: 1,
			Amount6:    big.NewInt(amountEach),
		}
	}
	return out
}

// BATCH-01：30 user → 切 2 批(25 + 5)。
func TestSettlement_BATCH01_Slicing(t *testing.T) {
	cli := &fakeSettlementClient{}
	store := newInMemBatchStore()
	svc := NewSettlementOrchestrator(cli, store)

	items := makeItems(30, 1_000) // 每条 1000，单批总额远低于上限
	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err)

	assert.Equal(t, 2, len(cli.calls), "30 user 应切 2 批")
	assert.Equal(t, 25, len(cli.calls[0].users), "第 1 批 25 个")
	assert.Equal(t, 5, len(cli.calls[1].users), "第 2 批 5 个")
	assert.Equal(t, 2, sum.ConfirmedBatches, "两批均成功确认")
	assert.Equal(t, 0, sum.FailedBatches)
	assert.Equal(t, MaxSettlementBatch, 25, "组批上限对齐 handoff §5.1 N≤25")
}

// BATCH-02：单批总额 > MAX_BATCH_TOTAL → 拒发/阻断（L2 绝对闸，复用 blockchain.MaxBatchTotal 同源常量）。
func TestSettlement_BATCH02_OverMaxBatchTotal(t *testing.T) {
	cli := &fakeSettlementClient{}
	store := newInMemBatchStore()
	svc := NewSettlementOrchestrator(cli, store)

	// 5 条，每条 MaxBatchTotal/4，使单批 sum > MaxBatchTotal。
	each := new(big.Int).Div(blockchain.MaxBatchTotal, big.NewInt(4))
	// each 可能 > MaxBillPerUser，但 L2 只测批总额绝对闸（不测单笔），用 each.Int64() 安全（MaxBatchTotal=1e10）。
	items := makeItems(5, each.Int64())

	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err, "L2 阻断不返回 fatal error，记 batch 状态供人工处理")

	assert.Equal(t, 0, len(cli.calls), "超 MAX_BATCH_TOTAL 绝对闸：拒发，不得上链")
	assert.Equal(t, 0, sum.ConfirmedBatches)
	assert.Equal(t, 1, sum.BlockedBatches, "该批被绝对闸阻断")
}

// CB-01：单批总额 > 历史月均 × N → 熔断 pending-review（不自动发）。
func TestSettlement_CB01_CircuitBreak(t *testing.T) {
	cli := &fakeSettlementClient{}
	store := newInMemBatchStore()
	// 历史月均 = (100 + 200)/2 = 150。N=CircuitBreakerMultiplier。阈值 = 150*N。
	store.historyTotals = []*big.Int{big.NewInt(100), big.NewInt(200)}
	svc := NewSettlementOrchestrator(cli, store)

	avg := int64(150)
	threshold := avg * int64(CircuitBreakerMultiplier)
	// 单批总额 = threshold + 一点，触发熔断（仍 < MaxBatchTotal 绝对闸，隔离两闸）。
	// 用 1 条，金额 = threshold+10。
	items := makeItems(1, threshold+10)

	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err)

	assert.Equal(t, 0, len(cli.calls), "超月均×N：熔断 pending-review，不自动发")
	assert.Equal(t, 1, sum.PendingReviewBatches, "该批进入人工放行")
	assert.Equal(t, 0, sum.ConfirmedBatches)

	// 落 DB 的批次状态应为 pending-review。
	rec, _ := store.GetBatch("2026-06", 0)
	require.NotNil(t, rec)
	assert.Equal(t, BatchStatusPendingReview, rec.Status)
}

// CB-02（红线）：冷启动无历史均值 → 不熔断，仅查绝对上限 MAX_BATCH_TOTAL，不得除零/失效。
func TestSettlement_CB02_ColdStartFallback(t *testing.T) {
	cli := &fakeSettlementClient{}
	store := newInMemBatchStore()
	store.historyTotals = nil // 冷启动：无历史样本（均值=0）
	svc := NewSettlementOrchestrator(cli, store)

	// 单批总额远超「月均×N」（因均值=0，若用月均闸会误熔断），但 < MAX_BATCH_TOTAL 绝对闸。
	// 应正常发出（仅绝对闸生效），证明冷启动回退、无除零失效。
	each := new(big.Int).Div(blockchain.MaxBatchTotal, big.NewInt(2)) // 单批 sum = MaxBatchTotal/2 < 绝对闸
	items := makeItems(1, each.Int64())

	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err, "冷启动不得因均值=0 失效/除零")

	assert.Equal(t, 1, len(cli.calls), "冷启动回退绝对闸：低于 MAX_BATCH_TOTAL 应正常发出，不误熔断")
	assert.Equal(t, 1, sum.ConfirmedBatches)
	assert.Equal(t, 0, sum.PendingReviewBatches, "冷启动无月均闸，不得误进 pending-review")
}

// IDEM-01：同 month+batchIndex 已 confirmed → 不重发。
func TestSettlement_IDEM01_ConfirmedNotResent(t *testing.T) {
	cli := &fakeSettlementClient{}
	store := newInMemBatchStore()
	svc := NewSettlementOrchestrator(cli, store)

	items := makeItems(30, 1_000) // 2 批

	// 第一次：两批都发出并确认。
	_, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cli.calls))

	// 第二次同 month：两批均已 confirmed → 不重发。
	sum2, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cli.calls), "已 confirmed 批不重发（month+batchIndex 幂等键）")
	assert.Equal(t, 2, sum2.SkippedBatches, "两批均因幂等跳过")
	assert.Equal(t, 0, sum2.ConfirmedBatches)
}

// FAIL-01：一批失败 → 其他批继续 + 失败批记录可重试。
func TestSettlement_FAIL01_FailedBatchRetriable(t *testing.T) {
	store := newInMemBatchStore()
	// 第 0 批失败（receipt status=0），第 1 批成功。
	cli := &fakeSettlementClient{
		resultFn: func(callIdx int, users []common.Address) (*gethtypes.Receipt, error) {
			if callIdx == 0 {
				return &gethtypes.Receipt{Status: gethtypes.ReceiptStatusFailed}, nil
			}
			return &gethtypes.Receipt{Status: gethtypes.ReceiptStatusSuccessful}, nil
		},
	}
	svc := NewSettlementOrchestrator(cli, store)

	items := makeItems(30, 1_000) // 2 批
	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err, "一批失败不应中断整体（失败批续跑）")

	assert.Equal(t, 2, len(cli.calls), "两批都尝试了——失败批不阻断其他批")
	assert.Equal(t, 1, sum.ConfirmedBatches, "第 1 批成功确认")
	assert.Equal(t, 1, sum.FailedBatches, "第 0 批失败记录")

	// 失败批落 DB 状态 failed，未被标记 confirmed → 可重试。
	rec0, _ := store.GetBatch("2026-06", 0)
	require.NotNil(t, rec0)
	assert.Equal(t, BatchStatusFailed, rec0.Status)
	rec1, _ := store.GetBatch("2026-06", 1)
	require.NotNil(t, rec1)
	assert.Equal(t, BatchStatusConfirmed, rec1.Status)

	// 重跑：失败批 (idx0) 重发并这次成功，已 confirmed 批 (idx1) 不重发。
	cli.resultFn = nil // 全部成功
	sum2, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err)
	assert.Equal(t, 3, len(cli.calls), "重跑只重发失败批（idx0），confirmed 批不重发")
	assert.Equal(t, 1, sum2.ConfirmedBatches, "失败批重试成功")
	assert.Equal(t, 1, sum2.SkippedBatches, "已 confirmed 批幂等跳过")
}
