package services

import (
	"math/big"
	"testing"

	"linkworld-backend/internal/blockchain"

	"github.com/stretchr/testify/assert"
)

// T5 计价引擎 TDD（design §4.2 / §7.0① L1 / arch-review B4）。
// 占位费率（PLACEHOLDER，design §4.2）：op=1 DataUnitPrice=10000(0.01 USDT/MB)、
// CallUnitPrice=5000(0.005 USDT/min)、MinBillAmount=0。单位锁：dataMB=MB、callMin=minute。

// PRICE-01：固定输入 → 精确断言 amount6（占位费率确定可复算）。
//   amount6 = 100*10000 + 10*5000 = 1_000_000 + 50_000 = 1_050_000（6 位最小单位 = 1.05 USDT）。
func TestPricing_PRICE01_FixedAmount(t *testing.T) {
	ps := NewPricingService()
	amount6, err := ps.Price(100, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1_050_000), amount6, "100MB*10000 + 10min*5000 = 1_050_000（6 位最小单位）")
}

// PRICE-02：usage 上界（B4 / L1）——dataMB 超 MaxDataMB 或 callMin 超 MaxCallMin → 拒绝 + error。
func TestPricing_PRICE02_UsageUpperBound(t *testing.T) {
	ps := NewPricingService()

	// dataMB 超界（MaxDataMB+1），callMin 合法 → 拒绝。
	_, err := ps.Price(MaxDataMB+1, 0, 1)
	assert.Error(t, err, "dataMB > MaxDataMB 必须拒绝")

	// callMin 超界（MaxCallMin+1），dataMB 合法 → 拒绝。
	_, err = ps.Price(0, MaxCallMin+1, 1)
	assert.Error(t, err, "callMin > MaxCallMin 必须拒绝")

	// 边界内（== 上限）不应因上界拒绝（可能因 L1 金额闸拒，单独在 PRICE-03 覆盖；
	// 这里用 0 费率算不超 amount 上限的安全输入：dataMB=MaxDataMB, callMin=0，op 用 0 费率 region）。
	// 为隔离上界与金额闸，这里仅断言「合法小输入不报上界错」。
	_, err = ps.Price(1, 1, 1)
	assert.NoError(t, err, "合法小输入不应被上界拒绝")
}

// PRICE-03：amount6 超单 bill 硬上限 MAX_BILL_PER_USER（L1）→ 拒绝；与 client L3 同常量（blockchain.MaxBillPerUser）。
func TestPricing_PRICE03_AmountOverMaxBillPerUser(t *testing.T) {
	ps := NewPricingService()

	// 构造一个 usage 使 amount6 > MaxBillPerUser，但 usage 仍在上界内。
	// op=1 DataUnitPrice=10000；MaxBillPerUser=1_000_000_000；
	// 需 dataMB*10000 > 1_000_000_000 → dataMB > 100_000，取 dataMB=200_000（< MaxDataMB=1_000_000，在上界内）。
	// amount6 = 200_000*10000 = 2_000_000_000 > MaxBillPerUser。
	assert.True(t, uint64(200_000) <= MaxDataMB, "前置：200_000 必须 ≤ MaxDataMB（隔离上界闸与金额闸）")
	_, err := ps.Price(200_000, 0, 1)
	assert.Error(t, err, "amount6 > MaxBillPerUser 必须拒绝（L1，与 L3 同常量）")

	// 锁定 L1 与 L3 用同一常量（单一事实源，防漂移）。
	assert.Equal(t, blockchain.MaxBillPerUser, MaxBillPerUser, "L1 与 L3 必须共用同一 MaxBillPerUser 常量")
}

// PRICE-04：纯整数无浮点 —— 不同输入精确值，table-driven 固定断言。
func TestPricing_PRICE04_PureIntegerNoFloat(t *testing.T) {
	ps := NewPricingService()
	cases := []struct {
		name    string
		dataMB  uint64
		callMin uint64
		op      uint
		want    *big.Int
	}{
		{"零用量出账0", 0, 0, 1, big.NewInt(0)},
		{"仅流量", 1, 0, 1, big.NewInt(10_000)},
		{"仅通话", 0, 1, 1, big.NewInt(5_000)},
		{"3MB+7min", 3, 7, 1, big.NewInt(3*10_000 + 7*5_000)},
		{"奇数无浮点截断", 999, 333, 1, big.NewInt(999*10_000 + 333*5_000)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ps.Price(tc.dataMB, tc.callMin, tc.op)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// PRICE-05：未知 operatorId 无费率 → 拒绝（防静默 0 出账打错运营商）。
func TestPricing_PRICE05_UnknownOperator(t *testing.T) {
	ps := NewPricingService()
	_, err := ps.Price(10, 10, 9999)
	assert.Error(t, err, "未知 operatorId 必须报错，不静默")
}
