package services

import (
	"fmt"
	"math/big"
	"testing"

	"linkworld-backend/internal/blockchain/bindings"
	"linkworld-backend/internal/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// OPID-01：seed Operator.ID == 链上 operatorId(1..11) 固定映射（design §4.5 / arch-review）。
// SeedOperators() 是 cmd/main.go seed 的单一事实源，显式写 ID=1..11，与 ServiceManager.initialize
// 注入顺序一一对应，不靠 name 比对。
func TestOperatorID_OPID01_FixedMapping(t *testing.T) {
	ops := SeedOperators()
	assert.Len(t, ops, 11, "11 个内置运营商")

	// 链上 ServiceManager.initialize 冻结顺序（id → countryCode），权威值。
	wantCountry := map[uint]string{
		1: "US", 2: "GB", 3: "FR", 4: "RU", 5: "JP", 6: "VN",
		7: "LA", 8: "KH", 9: "TH", 10: "MY", 11: "PH",
	}
	for i, op := range ops {
		wantID := uint(i + 1)
		assert.Equal(t, wantID, op.ID, "第 %d 个 seed 的 ID 必须 == 链上 operatorId %d（不靠 name 比对）", i+1, wantID)
		assert.Equal(t, wantCountry[wantID], op.CountryCode, "operatorId %d 的国家码必须与链上 initialize 顺序一致", wantID)
	}

	// OperatorID 映射函数：供 T6 组批用，后端 opID 即链上 operatorId（恒等，但显式提供入口）。
	assert.Equal(t, big.NewInt(7), ChainOperatorID(7), "ChainOperatorID 直接返回链上 operatorId（恒等映射）")
}

// fakeOperatorReader 是链上 ServiceManager 读取器的测试替身。
type fakeOperatorReader struct {
	byID map[uint]bindings.IServiceManagerOperator
	err  error
}

func (f *fakeOperatorReader) GetOperator(id *big.Int) (bindings.IServiceManagerOperator, error) {
	if f.err != nil {
		return bindings.IServiceManagerOperator{}, f.err
	}
	op, ok := f.byID[uint(id.Uint64())]
	if !ok {
		return bindings.IServiceManagerOperator{}, fmt.Errorf("operator %s 链上不存在", id)
	}
	return op, nil
}

// chainSeed 构造与 seed countryCode 对齐、paymentAddress 非零的链上数据（一致场景）。
func chainSeed(consistent bool) map[uint]bindings.IServiceManagerOperator {
	out := map[uint]bindings.IServiceManagerOperator{}
	nonZero := common.HexToAddress("0x000000000000000000000000000000000000dEaD")
	for _, op := range SeedOperators() {
		cc := op.CountryCode
		if !consistent && op.ID == 9 {
			cc = "XX" // 故意制造 operatorId=9 国家码漂移
		}
		out[op.ID] = bindings.IServiceManagerOperator{
			Id:             new(big.Int).SetUint64(uint64(op.ID)),
			Name:           op.Name,
			CountryCode:    cc,
			PaymentAddress: nonZero,
			IsActive:       true,
		}
	}
	return out
}

// OPID-02：sanity check 一致 → nil；不一致（国家码漂移）→ 返回不一致清单（由调用方决定 fail/warn）。
func TestOperatorID_OPID02_SanityCheck(t *testing.T) {
	seed := SeedOperators()

	// 一致：无 mismatch。
	reader := &fakeOperatorReader{byID: chainSeed(true)}
	mism, err := SanityCheckOperators(reader, seed)
	assert.NoError(t, err)
	assert.Empty(t, mism, "链上与 seed 一致 → 无 mismatch")

	// 不一致：operatorId=9 国家码漂移 → 出现在 mismatch 清单。
	reader2 := &fakeOperatorReader{byID: chainSeed(false)}
	mism2, err := SanityCheckOperators(reader2, seed)
	assert.NoError(t, err)
	assert.NotEmpty(t, mism2, "链上漂移 → 必须报告 mismatch（防分账打错地址）")
	assert.Contains(t, mism2[0], "9", "mismatch 应定位到 operatorId 9")
}

// OPID-03：链上 paymentAddress 为零地址 → 标记不可结算（mismatch），合约 createBill 也会 fail-fast。
func TestOperatorID_OPID03_ZeroPaymentAddress(t *testing.T) {
	chain := chainSeed(true)
	// operatorId=3 链上分账地址为零（未上链/未设置）。
	op3 := chain[3]
	op3.PaymentAddress = common.Address{}
	chain[3] = op3

	reader := &fakeOperatorReader{byID: chain}
	mism, err := SanityCheckOperators(reader, SeedOperators())
	assert.NoError(t, err)
	assert.NotEmpty(t, mism, "零分账地址 → 报告 mismatch（不可结算）")
	assert.Contains(t, mism[0], "3")
}

// 确保 SeedOperators 与 models.Operator 字段对齐（编译期 + 运行期双保险）。
func TestOperatorID_SeedShape(t *testing.T) {
	var _ []models.Operator = SeedOperators()
}
