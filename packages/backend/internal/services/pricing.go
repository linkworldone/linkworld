package services

import (
	"fmt"
	"log"
	"math/big"

	"linkworld-backend/internal/blockchain"
)

// ─────────────────────────────────────────────────────────────────────────────
// 计价引擎（design §4.2 / §7.0① 金额硬闸 L1 / arch-review B4 usage 上界）
//
// 单位锁（design §4.2 / arch-review）：
//   - dataMB：流量，单位 = MB（兆字节，整数）
//   - callMin：通话，单位 = minute（分钟，整数）
//   - 单价 / 金额：USDT 6 位最小单位（整数，与合约 uint256 + usdtDecimals=6 对齐）
//
// 全程纯 *big.Int 整数运算，杜绝浮点（design §4.2 设计要点①）。
//
// 金额硬闸三层中的 L1（计价层，design §7.0①）：
//   ① usage 上界（B4）：dataMB ≤ MaxDataMB && callMin ≤ MaxCallMin，超界拒绝 + 告警
//   ② 单 bill 硬上限：amount6 ≤ MaxBillPerUser（== blockchain.MaxBillPerUser，与 L3 同一常量）
// L2 组批熔断（T6）/ L3 client 发交易前断言（T3）是另两层，三闸独立。
// ─────────────────────────────────────────────────────────────────────────────

// usage 上界 PLACEHOLDER 占位常量（design §0.1 L6 / §7.0②）。
//
// !!! PLACEHOLDER：以下阈值非真实业务值，上线前必须由产品/运营/安全确认并挂真机验收 gate。
// 启动时 log.Warn 提示占位（见 init）。测试用这些固定值断言。!!!
const (
	// MaxDataMB PLACEHOLDER：单 (user,operator) 单月流量上界 = 1,000,000 MB。
	MaxDataMB uint64 = 1_000_000
	// MaxCallMin PLACEHOLDER：单 (user,operator) 单月通话上界 = 100,000 min。
	MaxCallMin uint64 = 100_000
)

// MaxBillPerUser 是 L1 单 bill 金额硬上限，直接复用 blockchain.MaxBillPerUser（单一事实源，
// 防 L1 与 L3 阈值漂移，design §7.0① 三层同语义）。6 位最小单位。
var MaxBillPerUser = blockchain.MaxBillPerUser

// OperatorRate 是单运营商的计价费率（design §4.2，费率值占位，结构锁定）。
// 全 *big.Int 6 位最小单位；OperatorID == models.Operator.ID == 链上 operatorId（§4.5 固定映射，不靠 name 比对）。
type OperatorRate struct {
	OperatorID    uint     // == 链上 operatorId（§4.5 固定映射）
	Region        string   // 冗余便于按地区批量配置
	DataUnitPrice *big.Int // 每 1 MB 流量单价（6 位最小单位 USDT）
	CallUnitPrice *big.Int // 每 1 分钟通话单价（6 位最小单位 USDT）
	MinBillAmount *big.Int // 最小出账金额（低于此不出账，避免 dust）
}

// PricingService 计价引擎：usage × 费率 → amount6（6 位最小单位），纯整数无浮点。
// 替换 OperatorAPISimulator.GetBill 的 rand.Intn 随机金额（design §4.2 / 用户拍板决策 1）。
//
// 注意（CEO 验收语言，arch-review）：本服务保证「计价引擎正确」（金额确定可复算），
// **不等于「计费正确」**——usage 本轮仍由 OperatorAPISimulator.GetUsage 模拟（真实运营商 API 留后续）。
type PricingService struct {
	rates map[uint]OperatorRate // key = operatorID（== 链上 operatorId）
}

// NewPricingService 用占位费率表构造计价引擎（design §4.2）。
//
// !!! PLACEHOLDER 费率，非真实定价，上线前必须由运营填实数并挂真机验收 gate；
// 启动时 log.Warn 提示占位（见 init / 本构造）。!!!
func NewPricingService() *PricingService {
	return &PricingService{rates: placeholderRates()}
}

// Price 计算 (dataMB, callMin) 在 operatorID 下的出账金额（6 位最小单位），纯 *big.Int 整数运算。
//
//	amount6 = dataMB*DataUnitPrice + callMin*CallUnitPrice
//
// 校验顺序（fail-fast，非静默）：
//  1. operatorID 必须有费率（未知 → 拒绝，防静默 0 出账 / 打错运营商）
//  2. usage 上界 L1（B4）：dataMB > MaxDataMB || callMin > MaxCallMin → 拒绝 + 告警
//  3. 单 bill 硬上限 L1：amount6 > MaxBillPerUser → 拒绝 + 告警（与 client L3 同一常量）
//  4. amount6 < MinBillAmount → 返回 0（不出账，避免 dust；amount==0 上游跳过）
func (s *PricingService) Price(dataMB, callMin uint64, operatorID uint) (*big.Int, error) {
	rate, ok := s.rates[operatorID]
	if !ok {
		return nil, fmt.Errorf("计价拒绝：operatorID=%d 无费率配置（不静默 0 出账，防打错运营商）", operatorID)
	}

	// ② usage 上界 L1（B4）：超界拒绝 + 告警，绝不进入计价。
	if dataMB > MaxDataMB {
		log.Printf("WARN: 计价拒绝(L1 上界)：dataMB=%d 超 MaxDataMB=%d（op=%d），防天文账单", dataMB, MaxDataMB, operatorID)
		return nil, fmt.Errorf("计价拒绝：dataMB=%d 超上界 MaxDataMB=%d", dataMB, MaxDataMB)
	}
	if callMin > MaxCallMin {
		log.Printf("WARN: 计价拒绝(L1 上界)：callMin=%d 超 MaxCallMin=%d（op=%d），防天文账单", callMin, MaxCallMin, operatorID)
		return nil, fmt.Errorf("计价拒绝：callMin=%d 超上界 MaxCallMin=%d", callMin, MaxCallMin)
	}

	// amount6 = dataMB*DataUnitPrice + callMin*CallUnitPrice（纯整数）。
	amount6 := new(big.Int).Mul(new(big.Int).SetUint64(dataMB), rate.DataUnitPrice)
	amount6.Add(amount6, new(big.Int).Mul(new(big.Int).SetUint64(callMin), rate.CallUnitPrice))

	// ③ 单 bill 硬上限 L1：amount6 > MaxBillPerUser → 拒绝 + 告警（与 L3 同常量）。
	if amount6.Cmp(MaxBillPerUser) > 0 {
		log.Printf("WARN: 计价拒绝(L1 金额闸)：amount6=%s 超 MaxBillPerUser=%s（op=%d）", amount6, MaxBillPerUser, operatorID)
		return nil, fmt.Errorf("计价拒绝：amount6=%s 超单 bill 上限 MaxBillPerUser=%s", amount6, MaxBillPerUser)
	}

	// ④ 低于最小出账额 → 0（不出账，合约侧 amount>0 才 createBill）。
	if rate.MinBillAmount != nil && amount6.Cmp(rate.MinBillAmount) < 0 {
		return big.NewInt(0), nil
	}
	return amount6, nil
}

// placeholderRates 返回 11 个内置运营商的占位费率（design §4.2，operatorID == 链上 1..11）。
//
// !!! PLACEHOLDER 费率值，非真实定价！上线前必须由运营填实数并挂真机验收 gate。!!!
// 占位口径：DataUnitPrice=10000(0.01 USDT/MB)、CallUnitPrice=5000(0.005 USDT/min)、MinBillAmount=0。
func placeholderRates() map[uint]OperatorRate {
	rates := make(map[uint]OperatorRate, 11)
	for _, op := range SeedOperators() {
		rates[op.ID] = OperatorRate{
			OperatorID:    op.ID,
			Region:        op.Region,
			DataUnitPrice: big.NewInt(10_000), // PLACEHOLDER 0.01 USDT/MB
			CallUnitPrice: big.NewInt(5_000),  // PLACEHOLDER 0.005 USDT/min
			MinBillAmount: big.NewInt(0),      // PLACEHOLDER 无 dust 门槛
		}
	}
	return rates
}

func init() {
	// CEO 要求：占位费率/上界刺眼化，防永久占位（design §0.1 注 / arch-review）。
	log.Printf("WARN: PRICING TABLE IS PLACEHOLDER —— 费率(0.01USDT/MB,0.005USDT/min)与 usage 上界(MaxDataMB=%d,MaxCallMin=%d) 均为占位值，上线前必须由产品/运营/安全填真实数并挂真机验收 gate", MaxDataMB, MaxCallMin)
}
