package services

import (
	"fmt"
	"math/big"

	"linkworld-backend/internal/blockchain/bindings"
	"linkworld-backend/internal/models"

	"github.com/ethereum/go-ethereum/common"
)

// ─────────────────────────────────────────────────────────────────────────────
// operatorId 固定映射（design §4.5 / arch-review 最高优先项 / plan-review OPID）
//
// 决策：后端 seed 显式写 Operator.ID = 链上 operatorId(1..11)，使 models.Operator.ID
// 成为「链上 operatorId」的单一事实源。monthlySettlement 传 operatorIds[] 时直接用
// models.Operator.ID，无需任何转换 / name 比对。
//
// 为什么不靠 name 比对：后端 seed "T-Mobile" 与链上 "T-Mobile US" 字符串不等会静默 miss
// → 分账打到错误运营商地址（最隐蔽资损）。固定 ID 是单一事实源。
//
// seed 顺序与 contracts/contracts/ServiceManager.sol initialize 的 _operators[1..11] 逐一对应
// （已核对 countryCode：US/GB/FR/RU/JP/VN/LA/KH/TH/MY/PH）。
// ─────────────────────────────────────────────────────────────────────────────

// SeedOperators 返回 11 个内置运营商，显式指定 ID=链上 operatorId(1..11)（GORM 允许显式主键）。
// 是 cmd/main.go seed 与 PricingService 费率表的单一事实源（避免两处 seed 漂移）。
//
// 注：name 用后端展示名（与链上 "T-Mobile US" 等略不同是有意的——映射靠 ID 不靠 name）。
// requiredDeposit 字符串语义见 models 注（本轮沿用现状，精度迁移非 T5 范围）。
func SeedOperators() []models.Operator {
	return []models.Operator{
		{ID: 1, Name: "T-Mobile", Region: "United States", CountryCode: "US", RequiredDeposit: "0.01", IsActive: true},
		{ID: 2, Name: "Vodafone", Region: "United Kingdom", CountryCode: "GB", RequiredDeposit: "0.008", IsActive: true},
		{ID: 3, Name: "Orange", Region: "France", CountryCode: "FR", RequiredDeposit: "0.008", IsActive: true},
		{ID: 4, Name: "MTS", Region: "Russia", CountryCode: "RU", RequiredDeposit: "0.005", IsActive: true},
		{ID: 5, Name: "SoftBank", Region: "Japan", CountryCode: "JP", RequiredDeposit: "0.012", IsActive: true},
		{ID: 6, Name: "Viettel", Region: "Vietnam", CountryCode: "VN", RequiredDeposit: "0.003", IsActive: true},
		{ID: 7, Name: "Unitel", Region: "Laos", CountryCode: "LA", RequiredDeposit: "0.003", IsActive: true},
		{ID: 8, Name: "Smart", Region: "Cambodia", CountryCode: "KH", RequiredDeposit: "0.003", IsActive: true},
		{ID: 9, Name: "AIS", Region: "Thailand", CountryCode: "TH", RequiredDeposit: "0.004", IsActive: true},
		{ID: 10, Name: "Maxis", Region: "Malaysia", CountryCode: "MY", RequiredDeposit: "0.004", IsActive: true},
		{ID: 11, Name: "Globe", Region: "Philippines", CountryCode: "PH", RequiredDeposit: "0.003", IsActive: true},
	}
}

// ChainOperatorID 把后端 operatorID 映射为链上 operatorId（big.Int）供 T6 组批传 monthlySettlement 用。
// 因 §4.5 固定映射，后端 ID 即链上 operatorId（恒等），但显式提供入口集中表达契约、便于后续演进。
func ChainOperatorID(operatorID uint) *big.Int {
	return new(big.Int).SetUint64(uint64(operatorID))
}

// OperatorChainReader 是链上 ServiceManager 运营商读取器（供 sanity check 注入测试替身）。
// bindings.ServiceManagerCaller 的 GetOperator 签名是 (opts, *big.Int)；这里收敛为按 ID 读，
// 由适配器（NewServiceManagerCallerReader）桥接，测试用 fake 实现。
type OperatorChainReader interface {
	GetOperator(operatorID *big.Int) (bindings.IServiceManagerOperator, error)
}

// SanityCheckOperators 启动期读链校验：seed 的每个 operatorId 在链上存在、countryCode 一致、
// paymentAddress 非零（可结算）。返回不一致清单（design §4.5：调用方据此 fail-fast 或降级 warn）。
//
//   - 不靠 name 比对（链上 "T-Mobile US" vs seed "T-Mobile"）：用 countryCode + 存在性 + 分账地址非零。
//   - error 仅在链读取本身失败时返回；业务不一致走 mismatch 清单（让 main.go 决定 fail/warn）。
//   - 零分账地址 → 标记 mismatch（不可结算，合约 createBill 亦 fail-fast，后端提前拦避免整批浪费 gas）。
func SanityCheckOperators(reader OperatorChainReader, seed []models.Operator) ([]string, error) {
	var mismatches []string
	for _, op := range seed {
		chainOp, err := reader.GetOperator(ChainOperatorID(op.ID))
		if err != nil {
			return nil, fmt.Errorf("读链 getOperator(%d) 失败: %w", op.ID, err)
		}
		// 存在性：链上 id 必须 == 请求 id（不存在时合约返回零值 struct，id==0）。
		if chainOp.Id == nil || chainOp.Id.Uint64() != uint64(op.ID) {
			mismatches = append(mismatches, fmt.Sprintf("operatorId=%d 链上不存在（返回 id=%v）", op.ID, chainOp.Id))
			continue
		}
		// countryCode 一致（不靠 name）。
		if chainOp.CountryCode != op.CountryCode {
			mismatches = append(mismatches, fmt.Sprintf("operatorId=%d countryCode 漂移：seed=%q 链上=%q", op.ID, op.CountryCode, chainOp.CountryCode))
		}
		// 分账地址非零（可结算）。
		if chainOp.PaymentAddress == (common.Address{}) {
			mismatches = append(mismatches, fmt.Sprintf("operatorId=%d 链上 paymentAddress 为零地址（不可结算，跳过其 bill）", op.ID))
		}
	}
	return mismatches, nil
}

// serviceManagerCallerReader 把 bindings.ServiceManagerCaller 适配为 OperatorChainReader
// （生产路径：main.go 读链 sanity check 用）。
type serviceManagerCallerReader struct {
	caller *bindings.ServiceManagerCaller
}

// NewServiceManagerCallerReader 构造生产用读取器。caller 为 nil（未上链/降级）返回 nil，
// 调用方据此跳过 sanity check（design §4.5：占位地址/未上链时降级跳过）。
func NewServiceManagerCallerReader(caller *bindings.ServiceManagerCaller) OperatorChainReader {
	if caller == nil {
		return nil
	}
	return &serviceManagerCallerReader{caller: caller}
}

func (r *serviceManagerCallerReader) GetOperator(operatorID *big.Int) (bindings.IServiceManagerOperator, error) {
	return r.caller.GetOperator(nil, operatorID)
}
