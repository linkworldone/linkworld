package blockchain

// T3 client 真实读/写 + owner 发交易 + nonce + WaitMined + 金额闸 L3 的 TDD 测试。
//
// 真实链交互：用 go-ethereum 的 simulated.Backend（真实 EVM + 真实区块/共识 faker）跑，
// 经 *真实的 abigen 绑定*（OracleTransactor/OracleCaller/DepositCaller）发交易/读链，断言
// 「编码 → 本地 nonce → 签名发交易 → WaitMined → receipt.Status / 解码」全链路——非 mock。
//
// ⚠️ 字节码取舍（见 testdata_bytecode_test.go 详注）：项目合约 cancun 编译（PUSH0/tstore），
// 而本仓库 pin 的 geth v1.13.5 simulated.Backend 受 ethash 共识所限只到 London，无法执行
// cancun 字节码。故此处部署 London 兼容的最小桩合约承接真实绑定的调用，验证 client 机制；
// 业务合约端到端结算属 T4/T8（hardhat 31337 / 升级 geth）。
//
// 纯逻辑断言（L3 金额闸 / owner key 缺失降级 / nonce 单调）不依赖任何合约执行。
//
// init 字节码生成片段：
//   python3 -c "rt=bytes.fromhex('00'); n=len(rt); print((bytes([0x60,n,0x60,0x0c,0x60,0,0x39,0x60,n,0x60,0,0xf3])+rt).hex())"

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// 测试用 owner 私钥（hardhat account#0 风格，仅测试，绝不用于真实链）。
const testOwnerPrivHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// simulated.Backend 的链 ID 固定为 1337。
const testChainID = uint64(1337)

// chainIDBackend 给 simulated.Backend 补一个 ChainID(ctx) 方法，满足 client 的
// chainIDProvider 接口（生产用的 *ethclient.Client 原生有此方法）。
type chainIDBackend struct {
	*backends.SimulatedBackend
	chainID *big.Int
}

func (b *chainIDBackend) ChainID(context.Context) (*big.Int, error) {
	return b.chainID, nil
}

// wrap 把 simulated.Backend 包成带 ChainID 的后端，chainID 取真实链值（1337）。
func wrap(b *backends.SimulatedBackend) *chainIDBackend {
	return &chainIDBackend{SimulatedBackend: b, chainID: new(big.Int).SetUint64(testChainID)}
}

// deployStub 部署一个最小桩合约（给定 runtime），返回地址。
func deployStub(t *testing.T, backend *backends.SimulatedBackend, auth *bind.TransactOpts, deployHex string) common.Address {
	t.Helper()
	code, err := hexutil.Decode(deployHex)
	if err != nil {
		t.Fatalf("decode stub bytecode: %v", err)
	}
	// 用空 ABI 部署纯字节码。
	addr, _, _, err := bind.DeployContract(auth, abi.ABI{}, code, backend)
	if err != nil {
		t.Fatalf("deploy stub: %v", err)
	}
	backend.Commit()
	return addr
}

// newTestEnv 启动 simulated.Backend，部署 STOP 桩（承接 Oracle 写）与 RET32 桩（承接读），
// 返回 backend、owner 地址、合约地址 map（键与 deployments.Proxies 一致："Oracle"/"Deposit"）。
func newTestEnv(t *testing.T) (*backends.SimulatedBackend, common.Address, map[string]common.Address) {
	t.Helper()
	key, err := crypto.HexToECDSA(testOwnerPrivHex)
	if err != nil {
		t.Fatalf("parse owner key: %v", err)
	}
	owner := crypto.PubkeyToAddress(key.PublicKey)

	alloc := core.GenesisAlloc{
		owner: {Balance: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))},
	}
	backend := backends.NewSimulatedBackend(alloc, 30_000_000)

	auth, err := bind.NewKeyedTransactorWithChainID(key, new(big.Int).SetUint64(testChainID))
	if err != nil {
		t.Fatalf("transactor: %v", err)
	}

	// Oracle 写交易 → STOP 桩（任何无返回值调用 status=1）。
	oracleAddr := deployStub(t, backend, auth, stopDeployBytecode)
	// Deposit 读调用 → RET32 桩（返回 32 零字节，uint256 解码为 0）。
	depositAddr := deployStub(t, backend, auth, ret32DeployBytecode)

	contracts := map[string]common.Address{
		"Oracle":  oracleAddr,
		"Deposit": depositAddr,
	}
	return backend, owner, contracts
}

// newReadEnv 让 Oracle 也指向 RET32 桩（VerifyServiceActive 走 OracleCaller 需要返回值）。
func newReadEnv(t *testing.T) (*backends.SimulatedBackend, common.Address, map[string]common.Address) {
	t.Helper()
	key, _ := crypto.HexToECDSA(testOwnerPrivHex)
	owner := crypto.PubkeyToAddress(key.PublicKey)
	alloc := core.GenesisAlloc{owner: {Balance: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))}}
	backend := backends.NewSimulatedBackend(alloc, 30_000_000)
	auth, _ := bind.NewKeyedTransactorWithChainID(key, new(big.Int).SetUint64(testChainID))

	ret32 := deployStub(t, backend, auth, ret32DeployBytecode)
	contracts := map[string]common.Address{"Oracle": ret32, "Deposit": ret32}
	return backend, owner, contracts
}

// autoCommit 在后台周期性 backend.Commit() 出块，让 client 内部的 bind.WaitMined 能等到
// 回执（simulated.Backend 不自动出块）。返回 stop 函数，调用后同步等待 goroutine 退出，
// 确保不会在 backend.Close() 之后再 Commit。真实链自动出块，无此问题。
func autoCommit(backend *backends.SimulatedBackend) (stop func()) {
	done := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		t := time.NewTicker(5 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-t.C:
				backend.Commit()
			}
		}
	}()
	return func() {
		close(done)
		<-exited // 等 goroutine 真正退出后再返回，避免与 backend.Close() 竞争
	}
}

// ─── 写：真实链交互（send → nonce → WaitMined → status） ──────────────────────

// TestMonthlySettlementSendsTxSuccessfully：owner 发 MonthlySettlement 交易，WaitMined status=1。
func TestMonthlySettlementSendsTxSuccessfully(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()

	c, err := NewClientWithBackend(wrap(backend), testChainID, contracts)
	if err != nil {
		t.Fatalf("NewClientWithBackend: %v", err)
	}
	if err := c.EnableOwnerWrites(testOwnerPrivHex); err != nil {
		t.Fatalf("EnableOwnerWrites: %v", err)
	}

	stop := autoCommit(backend)
	defer stop()
	// 一笔合法结算（amount 在 L3 区间内）。
	receipt, err := c.MonthlySettlement(context.Background(),
		[]common.Address{{0x1}}, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(1000)})
	if err != nil {
		t.Fatalf("MonthlySettlement: %v", err)
	}
	if receipt == nil || receipt.Status != 1 {
		t.Fatalf("expected receipt.Status=1, got %+v", receipt)
	}
}

// TestNonceNoCollisionAcrossBatches：连发多笔交易 nonce 不复用（本地计数器递增）。
func TestNonceNoCollisionAcrossBatches(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, _ := NewClientWithBackend(wrap(backend), testChainID, contracts)
	if err := c.EnableOwnerWrites(testOwnerPrivHex); err != nil {
		t.Fatalf("EnableOwnerWrites: %v", err)
	}

	stop := autoCommit(backend)
	defer stop()
	start := c.OwnerNonce() // 0，首次发交易时 resync 到链上 pending nonce
	for i := 0; i < 3; i++ {
		receipt, err := c.MonthlySettlement(context.Background(),
			[]common.Address{{0x1}}, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(1000)})
		if err != nil {
			t.Fatalf("settlement %d: %v", i, err)
		}
		if receipt.Status != 1 {
			t.Fatalf("settlement %d status=%d", i, receipt.Status)
		}
	}
	// 三笔交易后本地计数器应较首笔 resync 后推进 3（且全部 status=1 证明 nonce 无复用）。
	got := c.OwnerNonce()
	if got < start+3 {
		t.Fatalf("nonce counter did not advance by 3: start=%d got=%d", start, got)
	}
}

// ─── L3 金额硬闸（纯逻辑，发交易前断言，arch-review B1 最后防线） ──────────────────

// TestL3GateRejectsAmountOverMaxBillPerUser：单笔 amount > MAX_BILL_PER_USER 拒发。
func TestL3GateRejectsAmountOverMaxBillPerUser(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, _ := NewClientWithBackend(wrap(backend), testChainID, contracts)
	_ = c.EnableOwnerWrites(testOwnerPrivHex)

	over := new(big.Int).Add(MaxBillPerUser, big.NewInt(1))
	_, err := c.MonthlySettlement(context.Background(),
		[]common.Address{{0x1}}, []*big.Int{big.NewInt(1)}, []*big.Int{over})
	if err == nil {
		t.Fatal("expected L3 rejection for amount > MaxBillPerUser, got nil")
	}
}

// TestL3GateRejectsBatchTotalOverMax：sum(amounts) > MAX_BATCH_TOTAL 拒发。
func TestL3GateRejectsBatchTotalOverMax(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, _ := NewClientWithBackend(wrap(backend), testChainID, contracts)
	_ = c.EnableOwnerWrites(testOwnerPrivHex)

	// 每笔等于单笔上限（合法），凑足够多笔使合计超 MaxBatchTotal。
	perUser := new(big.Int).Set(MaxBillPerUser)
	n := new(big.Int).Div(MaxBatchTotal, MaxBillPerUser)
	count := int(n.Int64()) + 2 // 确保 sum > MaxBatchTotal
	users := make([]common.Address, count)
	opIds := make([]*big.Int, count)
	amounts := make([]*big.Int, count)
	for i := 0; i < count; i++ {
		users[i] = common.BigToAddress(big.NewInt(int64(i + 1)))
		opIds[i] = big.NewInt(1)
		amounts[i] = new(big.Int).Set(perUser)
	}
	_, err := c.MonthlySettlement(context.Background(), users, opIds, amounts)
	if err == nil {
		t.Fatal("expected L3 rejection for sum > MaxBatchTotal, got nil")
	}
}

// TestL3GateRejectsZeroAmount：amount=0 拒发（L3 区间 (0, MAX]）。
func TestL3GateRejectsZeroAmount(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, _ := NewClientWithBackend(wrap(backend), testChainID, contracts)
	_ = c.EnableOwnerWrites(testOwnerPrivHex)

	_, err := c.MonthlySettlement(context.Background(),
		[]common.Address{{0x1}}, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(0)})
	if err == nil {
		t.Fatal("expected L3 rejection for amount=0, got nil")
	}
}

// ─── owner key 缺失降级 + chainID 校验 ────────────────────────────────────────

// TestOwnerKeyMissingDegradesWrites：未注入 owner key → 写方法返回明确 error，不 panic。
func TestOwnerKeyMissingDegradesWrites(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, err := NewClientWithBackend(wrap(backend), testChainID, contracts)
	if err != nil {
		t.Fatalf("NewClientWithBackend: %v", err)
	}
	// 未调用 EnableOwnerWrites。
	if c.CanWrite() {
		t.Fatal("CanWrite should be false without owner key")
	}
	_, err = c.MonthlySettlement(context.Background(),
		[]common.Address{{0x1}}, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(1000)})
	if err == nil {
		t.Fatal("expected error when writes disabled, got nil")
	}
}

// TestChainIDMismatchRejected：注入的 chainID 与链上 ChainID 不一致 → EnableOwnerWrites 报错。
func TestChainIDMismatchRejected(t *testing.T) {
	backend, _, contracts := newTestEnv(t)
	defer backend.Close()
	c, err := NewClientWithBackend(wrap(backend), testChainID+999, contracts)
	if err != nil {
		t.Fatalf("NewClientWithBackend: %v", err)
	}
	if err := c.EnableOwnerWrites(testOwnerPrivHex); err == nil {
		t.Fatal("expected chainID mismatch error, got nil")
	}
	if c.CanWrite() {
		t.Fatal("writes must stay disabled after chainID mismatch")
	}
}

// ─── 读：真实链交互（通过真实 Caller 解码返回值） ──────────────────────────────

// TestReadMethodsReturnOnChainValues：读方法走真实 Caller 解码链上返回（桩返回 0/false）。
func TestReadMethodsReturnOnChainValues(t *testing.T) {
	backend, owner, contracts := newReadEnv(t)
	defer backend.Close()
	c, err := NewClientWithBackend(wrap(backend), testChainID, contracts)
	if err != nil {
		t.Fatalf("NewClientWithBackend: %v", err)
	}

	amt, err := c.GetDepositAmount(owner.Hex())
	if err != nil {
		t.Fatalf("GetDepositAmount: %v", err)
	}
	if amt.Sign() != 0 {
		t.Fatalf("expected 0 deposit, got %s", amt)
	}

	exp, err := c.GetLockExpiry(owner.Hex())
	if err != nil {
		t.Fatalf("GetLockExpiry: %v", err)
	}
	if exp != 0 {
		t.Fatalf("expected 0 lock expiry, got %d", exp)
	}

	active, err := c.VerifyServiceActive(owner.Hex())
	if err != nil {
		t.Fatalf("VerifyServiceActive: %v", err)
	}
	if active {
		t.Fatal("expected service inactive, got active")
	}
}
