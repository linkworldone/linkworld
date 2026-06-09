package services

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"linkworld-backend/internal/blockchain"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ─── T6 组批结算编排（design §6.1 + §7.0① L2 组批层 + arch-review B1-L2） ───────────
//
// 职责：把每 user 的 (链上 operatorId, amount6) 组装为三等长数组 → 切片每批 ≤25（handoff §5.1）
//       → 逐批过 L2 金额闸（绝对上限 + 历史月均×N 熔断，含冷启动回退红线）→ 调 client.MonthlySettlement
//       → 据 receipt.Status 判成败 → month+batchIndex 幂等键落 DB（已 confirmed 不重发，失败批可重试）。
//
// 三层金额硬闸定位（design §7.0①）：本编排是 L2（组批层）。L1 计价层（PricingService.Price，T5）/
// L3 client 发交易前断言（assertAmountGateL3，T3）独立存在，即使本层漏校验 L3 仍兜底。

const (
	// MaxSettlementBatch 单批 user 上限（handoff §5.1 monthlySettlement 全链路最紧约束 N≤25）。
	MaxSettlementBatch = 25

	// CircuitBreakerMultiplier L2 异常熔断倍数 N（PLACEHOLDER，design §7.0②/§11，占位 3）：
	// 单批总额 > 历史月均 × N → 熔断暂停，标记 pending_review 待人工放行（不自动发）。
	// 真值待产品/运营/安全确认。
	CircuitBreakerMultiplier = 3
)

// 批次状态常量别名（单一事实源在 models，编排层/测试复用）。
const (
	BatchStatusPending       = models.BatchStatusPending
	BatchStatusConfirmed     = models.BatchStatusConfirmed
	BatchStatusFailed        = models.BatchStatusFailed
	BatchStatusPendingReview = models.BatchStatusPendingReview
)

func init() {
	// 占位常量刺眼化（design §0.1/§10）。绝对闸 MaxBatchTotal 与 L3 同源（blockchain 包，单一事实源）。
	log.Printf("WARN: L2 组批熔断倍数 N=%d 为 PLACEHOLDER（design §7.0②），上线前须由产品/运营/安全确认；"+
		"绝对闸 MaxBatchTotal=%s 复用 L3 同源常量（blockchain 包）", CircuitBreakerMultiplier, blockchain.MaxBatchTotal)
}

// SettlementItem 是单个 user 的结算项（已过 L1 计价/上界）。Amount6 为 6 位最小单位。
type SettlementItem struct {
	User       common.Address
	OperatorID uint     // 后端 operatorID；发链前用 ChainOperatorID 映射
	Amount6    *big.Int // 6 位最小单位
}

// SettlementSummary 是一次 SettleMonth 的分批结果摘要（供 handler 返回 + 运维观测）。
type SettlementSummary struct {
	Month                string
	TotalBatches         int
	ConfirmedBatches     int      // 本次发出并 receipt.Status=1
	FailedBatches        int      // 本次发出但 receipt.Status=0 或发送失败（可重试）
	SkippedBatches       int      // 已 confirmed，幂等跳过未重发
	BlockedBatches       int      // 超 MAX_BATCH_TOTAL 绝对闸，拒发
	PendingReviewBatches int      // 超历史月均×N 熔断，待人工放行
	TxHashes             []string // confirmed 批的 txHash
}

// SettlementClient 抽象链上写调用，供测试 mock（生产由 *blockchain.Client 满足，签名一致）。
type SettlementClient interface {
	MonthlySettlement(ctx context.Context, users []common.Address, operatorIds, amounts []*big.Int) (*gethtypes.Receipt, error)
}

// 编译期断言：生产 *blockchain.Client 必须满足 SettlementClient（签名漂移即编译失败）。
var _ SettlementClient = (*blockchain.Client)(nil)

// settlementBatchRecord 是编排层视角的批次记录（与 models.SettlementBatch 解耦，便于 store 替身）。
type settlementBatchRecord struct {
	Month       string
	BatchIndex  int
	UserCount   int
	TotalAmount *big.Int
	Status      string
	TxHash      string
	Note        string
}

// SettlementBatchStore 抽象批次幂等存储（生产由 settlementBatchRepoStore 适配 repository，测试用内存替身）。
type SettlementBatchStore interface {
	GetBatch(month string, batchIndex int) (*settlementBatchRecord, error)
	SaveBatch(r *settlementBatchRecord) error
	// HistoricalBatchTotals 返回 beforeMonth 之前已 confirmed 批的总额清单（供月均熔断）。
	HistoricalBatchTotals(beforeMonth string) ([]*big.Int, error)
}

// SettlementOrchestrator 编排组批结算。
type SettlementOrchestrator struct {
	client SettlementClient
	store  SettlementBatchStore
}

func NewSettlementOrchestrator(client SettlementClient, store SettlementBatchStore) *SettlementOrchestrator {
	return &SettlementOrchestrator{client: client, store: store}
}

// SettleMonth 对一个月份的全部结算项组批并逐批结算（幂等可重跑）。
//   - 切片每批 ≤ MaxSettlementBatch。
//   - 逐批 L2 闸：① 绝对闸 sum ≤ MaxBatchTotal（超 → BlockedBatches，拒发）；
//     ② 月均熔断 sum > 月均×N（超 → PendingReviewBatches，不自动发）；
//     ③ 冷启动回退（红线）：无历史样本 → 跳过月均闸，仅查绝对闸，绝不除零/失效。
//   - month+batchIndex 幂等：已 confirmed 批跳过不重发；failed/pending/不存在 批发送。
//   - 失败批续跑：单批失败/阻断不影响其他批。
//
// 返回 error 仅用于不可恢复的存储/前置错误；单批业务失败记入 summary，不中断整体。
func (o *SettlementOrchestrator) SettleMonth(ctx context.Context, month string, items []SettlementItem) (*SettlementSummary, error) {
	sum := &SettlementSummary{Month: month}

	// 计算历史月均（冷启动回退红线：无样本时 avg=nil，表示月均闸不生效）。
	avg, err := o.historicalAverage(month)
	if err != nil {
		return nil, fmt.Errorf("读历史月均失败：%w", err)
	}

	batches := chunkItems(items, MaxSettlementBatch)
	sum.TotalBatches = len(batches)

	for idx, batch := range batches {
		if err := o.settleOneBatch(ctx, month, idx, batch, avg, sum); err != nil {
			// 存储类不可恢复错误：记 failed 并继续（不中断其他批，失败批续跑）。
			log.Printf("WARN: 结算批 month=%s idx=%d 存储错误：%v（记 failed，续跑其他批）", month, idx, err)
			sum.FailedBatches++
		}
	}
	return sum, nil
}

// settleOneBatch 处理单批：幂等检查 → L2 闸 → 发交易 → 据 receipt 落状态。
func (o *SettlementOrchestrator) settleOneBatch(
	ctx context.Context, month string, idx int, batch []SettlementItem, avg *big.Int, sum *SettlementSummary,
) error {
	// 1. 幂等：已 confirmed 批不重发（month+batchIndex 幂等键）。
	prev, err := o.store.GetBatch(month, idx)
	if err != nil {
		return err
	}
	if prev != nil && prev.Status == models.BatchStatusConfirmed {
		sum.SkippedBatches++
		if prev.TxHash != "" {
			sum.TxHashes = append(sum.TxHashes, prev.TxHash)
		}
		return nil
	}

	users, operatorIds, amounts, total := buildArrays(batch)
	rec := &settlementBatchRecord{
		Month:       month,
		BatchIndex:  idx,
		UserCount:   len(users),
		TotalAmount: total,
		Status:      models.BatchStatusPending,
	}

	// 2. L2 闸①：绝对上限（复用 blockchain.MaxBatchTotal，与 L3 同源单一常量）。
	if total.Cmp(blockchain.MaxBatchTotal) > 0 {
		rec.Status = models.BatchStatusPendingReview
		rec.Note = fmt.Sprintf("L2 绝对闸：批总额 %s 超 MaxBatchTotal=%s，拒发", total, blockchain.MaxBatchTotal)
		log.Printf("ALERT: %s（month=%s idx=%d）", rec.Note, month, idx)
		sum.BlockedBatches++
		return o.store.SaveBatch(rec)
	}

	// 3. L2 闸②：历史月均×N 熔断（冷启动回退红线：avg==nil → 跳过此闸，仅绝对闸生效，不除零）。
	if avg != nil && avg.Sign() > 0 {
		threshold := new(big.Int).Mul(avg, big.NewInt(int64(CircuitBreakerMultiplier)))
		if total.Cmp(threshold) > 0 {
			rec.Status = models.BatchStatusPendingReview
			rec.Note = fmt.Sprintf("L2 熔断：批总额 %s 超历史月均×N（%s×%d=%s），待人工放行",
				total, avg, CircuitBreakerMultiplier, threshold)
			log.Printf("ALERT: %s（month=%s idx=%d）", rec.Note, month, idx)
			sum.PendingReviewBatches++
			return o.store.SaveBatch(rec)
		}
	}

	// 4. 发交易（L3 client 发交易前还会再断言一次，最后防线）。
	receipt, err := o.client.MonthlySettlement(ctx, users, operatorIds, amounts)
	if err != nil {
		// 发送失败（含 L3 闸拒发/nonce 等）：记 failed 可重试（client 内部已 markNonceDirty）。
		rec.Status = models.BatchStatusFailed
		rec.Note = fmt.Sprintf("发送失败：%v", err)
		log.Printf("WARN: 结算批发送失败 month=%s idx=%d：%v（记 failed 可重试）", month, idx, err)
		sum.FailedBatches++
		return o.store.SaveBatch(rec)
	}

	if receipt != nil && receipt.Status == gethtypes.ReceiptStatusSuccessful {
		rec.Status = models.BatchStatusConfirmed
		if receipt.TxHash != (common.Hash{}) {
			rec.TxHash = receipt.TxHash.Hex()
			sum.TxHashes = append(sum.TxHashes, rec.TxHash)
		}
		sum.ConfirmedBatches++
	} else {
		// receipt.Status=0（链上 revert）：记 failed 可重试。
		rec.Status = models.BatchStatusFailed
		rec.Note = "receipt.Status=0（链上执行失败）"
		if receipt != nil && receipt.TxHash != (common.Hash{}) {
			rec.TxHash = receipt.TxHash.Hex()
		}
		log.Printf("WARN: 结算批回执失败 month=%s idx=%d（receipt.Status=0，记 failed 可重试）", month, idx)
		sum.FailedBatches++
	}
	return o.store.SaveBatch(rec)
}

// historicalAverage 计算历史月均（冷启动回退红线：无样本返回 nil，表示月均闸不生效）。
func (o *SettlementOrchestrator) historicalAverage(month string) (*big.Int, error) {
	totals, err := o.store.HistoricalBatchTotals(month)
	if err != nil {
		return nil, err
	}
	if len(totals) == 0 {
		return nil, nil // 冷启动/无样本：不除零，回退绝对闸
	}
	sum := new(big.Int)
	for _, t := range totals {
		if t != nil {
			sum.Add(sum, t)
		}
	}
	return new(big.Int).Div(sum, big.NewInt(int64(len(totals)))), nil
}

// chunkItems 把结算项切片为每批 ≤ size。
func chunkItems(items []SettlementItem, size int) [][]SettlementItem {
	if size <= 0 {
		size = MaxSettlementBatch
	}
	var out [][]SettlementItem
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		out = append(out, items[i:end])
	}
	return out
}

// buildArrays 构造三等长数组（users/operatorIds/amounts）+ 批总额。operatorId 用 ChainOperatorID 映射。
func buildArrays(batch []SettlementItem) (users []common.Address, operatorIds, amounts []*big.Int, total *big.Int) {
	users = make([]common.Address, len(batch))
	operatorIds = make([]*big.Int, len(batch))
	amounts = make([]*big.Int, len(batch))
	total = new(big.Int)
	for i, it := range batch {
		users[i] = it.User
		operatorIds[i] = ChainOperatorID(it.OperatorID)
		amounts[i] = it.Amount6
		total.Add(total, it.Amount6)
	}
	return users, operatorIds, amounts, total
}

// ─── 生产适配：把 repository.SettlementBatchRepository 适配为 SettlementBatchStore ────────

type settlementBatchRepoStore struct {
	repo *repository.SettlementBatchRepository
}

// NewSettlementBatchRepoStore 包装 repo 为编排层 store（DB 持久幂等 + 历史月均）。
// 自建表（不依赖 main.go AutoMigrate 清单，与 event_sync 同策略）；建表失败 fail-fast。
func NewSettlementBatchRepoStore(repo *repository.SettlementBatchRepository) (SettlementBatchStore, error) {
	if err := repo.Migrate(); err != nil {
		return nil, fmt.Errorf("settlement_batches 建表失败：%w", err)
	}
	return &settlementBatchRepoStore{repo: repo}, nil
}

func (s *settlementBatchRepoStore) GetBatch(month string, batchIndex int) (*settlementBatchRecord, error) {
	m, err := s.repo.Get(month, batchIndex)
	if err != nil || m == nil {
		return nil, err
	}
	total, ok := new(big.Int).SetString(m.TotalAmount, 10)
	if !ok {
		// string→big.Int 校验 fail-fast 不静默（design §6.3 收紧项）。
		return nil, fmt.Errorf("批次 TotalAmount 非法整数：%q（month=%s idx=%d）", m.TotalAmount, month, batchIndex)
	}
	return &settlementBatchRecord{
		Month:       m.Month,
		BatchIndex:  m.BatchIndex,
		UserCount:   m.UserCount,
		TotalAmount: total,
		Status:      m.Status,
		TxHash:      m.TxHash,
		Note:        m.Note,
	}, nil
}

func (s *settlementBatchRepoStore) SaveBatch(r *settlementBatchRecord) error {
	totalStr := "0"
	if r.TotalAmount != nil {
		totalStr = r.TotalAmount.String()
	}
	return s.repo.Save(&models.SettlementBatch{
		Month:       r.Month,
		BatchIndex:  r.BatchIndex,
		UserCount:   r.UserCount,
		TotalAmount: totalStr,
		Status:      r.Status,
		TxHash:      r.TxHash,
		Note:        r.Note,
	})
}

func (s *settlementBatchRepoStore) HistoricalBatchTotals(beforeMonth string) ([]*big.Int, error) {
	strs, err := s.repo.HistoricalConfirmedTotals(beforeMonth)
	if err != nil {
		return nil, err
	}
	out := make([]*big.Int, 0, len(strs))
	for _, str := range strs {
		v, ok := new(big.Int).SetString(str, 10)
		if !ok {
			return nil, fmt.Errorf("历史批次 TotalAmount 非法整数：%q", str)
		}
		out = append(out, v)
	}
	return out, nil
}
