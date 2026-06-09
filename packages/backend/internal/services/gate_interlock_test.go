package services

// T8 三层金额闸联动（design §7.0① + arch-review B1 + §8「金额闸不漂移」）。
//
// 三层闸（design §7.0① 表）：
//   L1 计价层 PricingService.Price        —— amount6 ≤ MaxBillPerUser
//   L2 组批层 SettlementOrchestrator      —— 单批 sum ≤ MaxBatchTotal（绝对闸）
//   L3 client 层 blockchain.Client        —— 每 amount∈(0,MaxBillPerUser]、sum≤MaxBatchTotal
//
// 单元层面三闸各自已覆盖（pricing_test / settlement_test / client_test）。本文件补「联动」缺口：
//   ① 三层共用同一常量（单一事实源，不漂移）——常量恒等断言；
//   ② 同一金额穿三层行为一致——刚好等于上限放行、超一点三层均拦，金额语义不在层间漂移。

import (
	"context"
	"math/big"
	"testing"

	"linkworld-backend/internal/blockchain"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GATE-LINK-01：三层共用同一金额常量（单一事实源，防漂移）。
// L1(services.MaxBillPerUser) 与 L3(blockchain.MaxBillPerUser) 同源；
// L2(blockchain.MaxBatchTotal) 与 L3 同源。任一层 fork 自己的常量都会让本测失败。
func TestGateInterlock_LINK01_SingleSourceConstants(t *testing.T) {
	// L1 == L3 单 bill 上限。
	assert.Same(t, blockchain.MaxBillPerUser, MaxBillPerUser,
		"L1 与 L3 必须是同一 *big.Int 指针（services.MaxBillPerUser = blockchain.MaxBillPerUser），不得各自定义")
	assert.Equal(t, 0, MaxBillPerUser.Cmp(blockchain.MaxBillPerUser),
		"L1/L3 单 bill 上限值必须相等")
	// L2 批总额上限取自 blockchain.MaxBatchTotal（与 L3 同源），settlement.go 直接引用该常量。
	assert.True(t, blockchain.MaxBatchTotal.Cmp(blockchain.MaxBillPerUser) > 0,
		"前置：批总额上限应大于单 bill 上限（否则单条就撑满批）")
}

// GATE-LINK-02：单 bill 金额穿 L1——恰等上限放行、超一点拒绝（边界一致）。
func TestGateInterlock_LINK02_L1Boundary(t *testing.T) {
	ps := NewPricingService()

	// 用 op=1 的 DataUnitPrice=10000 反推「恰好等于 MaxBillPerUser」的 dataMB。
	// MaxBillPerUser=1_000_000_000；dataMB = 1_000_000_000/10000 = 100_000（≤ MaxDataMB=1_000_000 在上界内）。
	atCap := new(big.Int).Div(MaxBillPerUser, big.NewInt(10000))
	require.True(t, atCap.IsUint64())
	dataAtCap := atCap.Uint64()

	// 恰好等于上限：L1 放行（amount6 == MaxBillPerUser，闸是 `>` 不是 `>=`）。
	amount6, err := ps.Price(dataAtCap, 0, 1)
	require.NoError(t, err, "金额恰等于 MaxBillPerUser 应放行（闸为 amount6 > 上限 才拒）")
	assert.Equal(t, 0, amount6.Cmp(MaxBillPerUser), "恰好等于上限")

	// 超一点（多 1 MB）：L1 拒绝。
	_, err = ps.Price(dataAtCap+1, 0, 1)
	assert.Error(t, err, "超 MaxBillPerUser 一点点 L1 必须拒绝")
}

// GATE-LINK-03：批总额穿 L2——恰等上限放行、超一点阻断（与 L1 单 bill 闸独立、金额语义一致）。
func TestGateInterlock_LINK03_L2Boundary(t *testing.T) {
	// 单条金额取 MaxBillPerUser（L1 恰好放行的最大单 bill），用足量条数把批总额顶到 MaxBatchTotal 边界。
	// MaxBatchTotal=10_000_000_000，MaxBillPerUser=1_000_000_000 → 10 条恰好撑满批。
	perBill := new(big.Int).Set(MaxBillPerUser)
	n := new(big.Int).Div(blockchain.MaxBatchTotal, perBill)
	require.True(t, n.IsInt64())
	cnt := int(n.Int64())
	require.LessOrEqual(t, cnt, MaxSettlementBatch, "前置：撑满批的条数应 ≤ 单批上限，否则会被切片，干扰 L2 边界测")

	// 恰好等于 MaxBatchTotal（cnt 条 × perBill）：L2 放行（发交易）。
	{
		cli := &fakeSettlementClient{}
		svc := NewSettlementOrchestrator(cli, newInMemBatchStore())
		items := makeBigItems(cnt, perBill)
		sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
		require.NoError(t, err)
		assert.Equal(t, 1, len(cli.calls), "批总额恰等 MaxBatchTotal 应放行发出（闸为 sum > 上限 才拒）")
		assert.Equal(t, 0, sum.BlockedBatches)
		assert.Equal(t, 1, sum.ConfirmedBatches)
	}

	// 超一点（最后一条 +1）：L2 绝对闸阻断（不发交易）。
	{
		cli := &fakeSettlementClient{}
		svc := NewSettlementOrchestrator(cli, newInMemBatchStore())
		items := makeBigItems(cnt, perBill)
		items[len(items)-1].Amount6 = new(big.Int).Add(perBill, big.NewInt(1))
		sum, err := svc.SettleMonth(context.Background(), "2026-07", items)
		require.NoError(t, err)
		assert.Equal(t, 0, len(cli.calls), "批总额超 MaxBatchTotal 一点 L2 必须阻断，不上链")
		assert.Equal(t, 1, sum.BlockedBatches)
	}
}

// GATE-LINK-04：L3 是发交易前最后防线——即使上游（L1/L2）漏校验，client 仍按同常量拦。
// 用 SettlementOrchestrator 直接喂一条「超 MaxBillPerUser 的 amount」给真实 client.assertAmountGateL3
// （经 SettlementClient 接口）。这里用 fake client 模拟「L3 拒发返回 error → 记 failed 可重试」的契约，
// 真实 L3 断言逻辑已在 blockchain/client_test 覆盖；本测锁定「L2 放行的单条若超单 bill 上限，
// 编排层据 client error 记 failed 不静默吞」——即 L3 拒发被编排层正确感知。
func TestGateInterlock_LINK04_L3LastLineFeedback(t *testing.T) {
	// fake client 模拟 L3：单条 amount > MaxBillPerUser → 返回 error（拒发）。
	cli := &fakeSettlementClient{
		resultFn: func(callIdx int, users []common.Address) (*gethtypes.Receipt, error) {
			return nil, assertL3Error{}
		},
	}
	svc := NewSettlementOrchestrator(cli, newInMemBatchStore())

	// 单条 amount = MaxBillPerUser（不触发 L2 绝对闸，单批 sum < MaxBatchTotal），
	// 但模拟 client L3 拒发 → 编排层应记 failed（可重试），不计 confirmed。
	items := makeBigItems(1, new(big.Int).Set(MaxBillPerUser))
	sum, err := svc.SettleMonth(context.Background(), "2026-06", items)
	require.NoError(t, err, "L3 拒发是单批失败，不应中断整体（失败批续跑）")
	assert.Equal(t, 1, sum.FailedBatches, "L3 拒发应被编排层记为 failed（可重试），不静默吞")
	assert.Equal(t, 0, sum.ConfirmedBatches, "L3 拒发的批绝不能被记为 confirmed")
}

type assertL3Error struct{}

func (assertL3Error) Error() string { return "L3 金额闸拒发（模拟）" }

// makeBigItems 造 n 条结算项，每项金额 amountEach（*big.Int，支持超 int64 的大额）。
func makeBigItems(n int, amountEach *big.Int) []SettlementItem {
	users := makeUsers(n)
	out := make([]SettlementItem, n)
	for i := 0; i < n; i++ {
		out[i] = SettlementItem{
			User:       users[i],
			OperatorID: 1,
			Amount6:    new(big.Int).Set(amountEach),
		}
	}
	return out
}
