package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"linkworld-backend/internal/blockchain/bindings"
)

// ─────────────────────────────────────────────────────────────────────────────
// 金额硬闸 L3 占位常量（design §7.0① / arch-review B1 最后防线）
//
// !!! PLACEHOLDER：以下阈值非真实业务值，上线前必须由产品/运营/安全确认并挂真机验收
// gate。启动时 log.Warn 提示占位。测试用这些固定值断言。!!!
//
// 语义：USDT 6 位最小单位（design §6.3，精度从 deployments.usdtDecimals 读，此处常量
// 与 6 位精度对齐，仅作发交易前的绝对上限断言）。
//   MaxBillPerUser  = 单笔 bill 金额上限，amount ∈ (0, MaxBillPerUser]
//   MaxBatchTotal   = 单批 sum(amounts) 上限
// L3 是三层金额硬闸的最后一道（L1 计价层 T5 / L2 组批熔断 T6 / L3 client 发交易前 = 本闸）。
// 即使 L1/L2 漏校验，L3 在发交易前仍兜底拒发。
// ─────────────────────────────────────────────────────────────────────────────
var (
	// MaxBillPerUser PLACEHOLDER：1,000 USDT = 1_000_000_000（6 位）。
	MaxBillPerUser = new(big.Int).SetUint64(1_000_000_000)
	// MaxBatchTotal PLACEHOLDER：10,000 USDT = 10_000_000_000（6 位）。
	MaxBatchTotal = new(big.Int).SetUint64(10_000_000_000)
)

func init() {
	log.Printf("WARN: 金额硬闸 L3 阈值为 PLACEHOLDER（MaxBillPerUser=%s, MaxBatchTotal=%s，6 位最小单位），上线前必须由产品/运营/安全确认真实值", MaxBillPerUser, MaxBatchTotal)
}

// contractBackend 是 client 做读/写/部署所需的 go-ethereum 后端抽象。
// *ethclient.Client 与 simulated.Backend 均满足此接口，便于离线 TDD。
type contractBackend interface {
	bind.ContractBackend
	bind.DeployBackend
}

// Client handles Ethereum blockchain interactions.
//
// 读：通过 abigen 绑定的 Caller 真实读链（design §6.1）。
// 写：owner 私钥经 NewKeyedTransactorWithChainID 发交易（design §6.2），本地 nonce 计数器
// 防复用，WaitMined 等回执判成败。owner key 仅内存，不落盘/不日志（仅打 owner address）。
// 缺 owner key → 写降级关闭，只读仍可（design §6.2 / arch-review owner=root）。
type Client struct {
	chainID   uint64
	contracts map[string]common.Address
	backend   contractBackend
	ethClient *ethclient.Client // 仅当通过 RPC 构造时非 nil（供 EthClient() 给 event_sync 用）

	// owner 写交易状态（缺 key 时全部零值，CanWrite()==false）。
	ownerKey  *ecdsa.PrivateKey
	ownerAddr common.Address

	// 本地 nonce 计数器（mutex 保护）。init 用 PendingNonceAt，发一笔 +1，出错 resync。
	nonceMu   sync.Mutex
	nonce     uint64
	nonceInit bool
}

// NewClient 通过 RPC 连接构造 Client（生产路径）。owner 写能力需另调 EnableOwnerWrites。
func NewClient(rpcURL string, chainID uint64, contracts map[string]common.Address) (*Client, error) {
	ec, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	return &Client{
		chainID:   chainID,
		contracts: contracts,
		backend:   ec,
		ethClient: ec,
	}, nil
}

// NewClientWithBackend 用注入的后端构造 Client（测试路径：simulated.Backend）。
// ethClient 保持 nil（EthClient() 返回 nil，event_sync 不在此路径启动）。
func NewClientWithBackend(backend contractBackend, chainID uint64, contracts map[string]common.Address) (*Client, error) {
	if backend == nil {
		return nil, fmt.Errorf("backend is nil")
	}
	return &Client{
		chainID:   chainID,
		contracts: contracts,
		backend:   backend,
	}, nil
}

// EthClient 返回底层 *ethclient.Client（仅 RPC 构造路径非 nil），供 event_sync(T4) 使用。
func (c *Client) EthClient() *ethclient.Client {
	return c.ethClient
}

// EnableOwnerWrites 注入 owner 私钥开启链上写能力（design §6.2 / §7.1）。
//   - 私钥仅存进程内存（c.ownerKey），绝不落盘、不入库、不进日志（仅打 owner address）。
//   - chainID 一致校验：client 的 chainID 必须 == 链上 ChainID()（防 replay / 签错链）。
//   - 校验通过后 init 本地 nonce 计数器。
//
// 缺 key（调用方不调本方法）→ 写降级关闭，CanWrite()==false，读仍可。
func (c *Client) EnableOwnerWrites(privHex string) error {
	key, err := crypto.HexToECDSA(privHex)
	if err != nil {
		return fmt.Errorf("invalid owner private key: %w", err) // 不打印 key 内容
	}

	// chainID 一致校验（design §6.5 / arch-review）：拒绝连 A 链发 B 链。
	onChainID, err := c.backendChainID(context.Background())
	if err != nil {
		return fmt.Errorf("query on-chain chainID: %w", err)
	}
	if onChainID.Uint64() != c.chainID {
		return fmt.Errorf("chainID 不一致：client=%d 链上=%s（拒绝开启写，防签错链/replay）", c.chainID, onChainID)
	}

	c.ownerKey = key
	c.ownerAddr = crypto.PubkeyToAddress(key.PublicKey)
	c.nonceInit = false // 下次发交易前 resync
	log.Printf("链上写已开启，owner=%s（私钥仅内存，不落盘）", c.ownerAddr.Hex())
	return nil
}

// CanWrite 报告是否已注入 owner key（写能力是否开启）。
func (c *Client) CanWrite() bool {
	return c.ownerKey != nil
}

// chainIDProvider 由能直接报告链上 chainID 的后端实现（*ethclient.Client 满足）。
// simulated.Backend 不直接实现，测试侧用适配器补齐（client_test.go）。
type chainIDProvider interface {
	ChainID(context.Context) (*big.Int, error)
}

// backendChainID 查询链上 chainID（用于 chainID 一致校验，防签错链/replay）。
func (c *Client) backendChainID(ctx context.Context) (*big.Int, error) {
	if p, ok := c.backend.(chainIDProvider); ok {
		return p.ChainID(ctx)
	}
	return nil, fmt.Errorf("backend 不支持 chainID 查询")
}

// OwnerNonce 返回当前本地 nonce 计数器值（供测试断言 nonce 单调不复用）。
func (c *Client) OwnerNonce() uint64 {
	c.nonceMu.Lock()
	defer c.nonceMu.Unlock()
	return c.nonce
}

// ─── 读：真实链查询（design §6.1，通过 abigen Caller） ──────────────────────────

// GetDepositAmount 读 Deposit 合约用户保证金余额（6 位最小单位）。
func (c *Client) GetDepositAmount(userAddress string) (*big.Int, error) {
	dep, err := c.depositCaller()
	if err != nil {
		return nil, err
	}
	return dep.GetDepositAmount(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(userAddress))
}

// GetLockExpiry 读 Deposit 合约用户锁仓到期时间（Unix 秒）。
func (c *Client) GetLockExpiry(userAddress string) (uint64, error) {
	dep, err := c.depositCaller()
	if err != nil {
		return 0, err
	}
	exp, err := dep.GetLockExpiry(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(userAddress))
	if err != nil {
		return 0, err
	}
	return exp.Uint64(), nil
}

// VerifyServiceActive 通过 Oracle 合约校验用户服务是否活跃。
func (c *Client) VerifyServiceActive(userAddress string) (bool, error) {
	addr, ok := c.contracts["Oracle"]
	if !ok {
		return false, fmt.Errorf("Oracle 合约地址未配置")
	}
	oc, err := bindings.NewOracleCaller(addr, c.backend)
	if err != nil {
		return false, err
	}
	return oc.VerifyServiceActive(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(userAddress))
}

func (c *Client) depositCaller() (*bindings.DepositCaller, error) {
	addr, ok := c.contracts["Deposit"]
	if !ok {
		return nil, fmt.Errorf("Deposit 合约地址未配置")
	}
	return bindings.NewDepositCaller(addr, c.backend)
}

// ─── 写：owner 发交易（design §6.2 + 金额闸 L3 + nonce + WaitMined） ──────────────

// MonthlySettlement 调 Oracle.monthlySettlement(users, operatorIds, amounts) 发交易。
//   - 缺 owner key → 写降级关闭，返回明确 error（不 panic）。
//   - L3 金额硬闸（发交易前断言，arch-review B1 最后防线）：每笔 amount ∈ (0, MaxBillPerUser]，
//     sum(amounts) ≤ MaxBatchTotal；越界拒发 + error，绝不上链。
//   - 本地 nonce 计数器取 nonce（发一笔 +1），失败时下次 resync。
//   - bind.WaitMined 等回执，返回 *types.Receipt（status=1 成功，0 失败；供 T6 失败批续跑用）。
//
// 注：月卡发放逻辑已从合约移除（Oracle.monthlySettlement 不再调 deposit.issueMonthlyTrafficCards），
// 后端亦不再提供对应发卡方法。
func (c *Client) MonthlySettlement(ctx context.Context, users []common.Address, operatorIds, amounts []*big.Int) (*types.Receipt, error) {
	if !c.CanWrite() {
		return nil, fmt.Errorf("链上写已降级关闭：owner key 未注入（只读仍可）")
	}
	if len(users) != len(operatorIds) || len(users) != len(amounts) {
		return nil, fmt.Errorf("数组长度不一致：users=%d operatorIds=%d amounts=%d", len(users), len(operatorIds), len(amounts))
	}
	// L3 金额硬闸（发交易前断言）。
	if err := assertAmountGateL3(amounts); err != nil {
		return nil, err
	}

	addr, ok := c.contracts["Oracle"]
	if !ok {
		return nil, fmt.Errorf("Oracle 合约地址未配置")
	}
	tx, err := bindings.NewOracleTransactor(addr, c.backend)
	if err != nil {
		return nil, err
	}
	opts, err := c.newTransactOpts(ctx)
	if err != nil {
		return nil, err
	}
	sentTx, err := tx.MonthlySettlement(opts, users, operatorIds, amounts)
	if err != nil {
		c.markNonceDirty() // 发送失败：下次 resync，避免计数器跑偏
		return nil, fmt.Errorf("send monthlySettlement: %w", err)
	}
	return c.waitMined(ctx, sentTx)
}

// assertAmountGateL3 实现 L3 金额硬闸：每笔 amount ∈ (0, MaxBillPerUser]，sum ≤ MaxBatchTotal。
func assertAmountGateL3(amounts []*big.Int) error {
	sum := new(big.Int)
	for i, a := range amounts {
		if a == nil {
			return fmt.Errorf("L3 闸：amounts[%d] 为 nil", i)
		}
		if a.Sign() <= 0 {
			return fmt.Errorf("L3 闸：amounts[%d]=%s 非正数（须 ∈ (0, %s]），拒发", i, a, MaxBillPerUser)
		}
		if a.Cmp(MaxBillPerUser) > 0 {
			return fmt.Errorf("L3 闸：amounts[%d]=%s 超单笔上限 MaxBillPerUser=%s，拒发", i, a, MaxBillPerUser)
		}
		sum.Add(sum, a)
	}
	if sum.Cmp(MaxBatchTotal) > 0 {
		return fmt.Errorf("L3 闸：批次合计 %s 超 MaxBatchTotal=%s，拒发", sum, MaxBatchTotal)
	}
	return nil
}

// newTransactOpts 构造带本地 nonce + signer 的 TransactOpts（chainID 已在 EnableOwnerWrites 校验）。
func (c *Client) newTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(c.ownerKey, new(big.Int).SetUint64(c.chainID))
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	n, err := c.takeNonce(ctx)
	if err != nil {
		return nil, err
	}
	opts.Nonce = new(big.Int).SetUint64(n)
	return opts, nil
}

// takeNonce 取下一个 nonce（本地计数器，mutex 保护）。首次/脏标记时用 PendingNonceAt resync。
func (c *Client) takeNonce(ctx context.Context) (uint64, error) {
	c.nonceMu.Lock()
	defer c.nonceMu.Unlock()
	if !c.nonceInit {
		n, err := c.backend.PendingNonceAt(ctx, c.ownerAddr)
		if err != nil {
			return 0, fmt.Errorf("resync nonce: %w", err)
		}
		c.nonce = n
		c.nonceInit = true
	}
	n := c.nonce
	c.nonce++
	return n, nil
}

// markNonceDirty 标记 nonce 计数器需在下次发交易前 resync（发送失败后调用）。
func (c *Client) markNonceDirty() {
	c.nonceMu.Lock()
	c.nonceInit = false
	c.nonceMu.Unlock()
}

// waitMined 等交易回执（design §6.1）。返回 receipt，据 receipt.Status 判成败（供 T6 用）。
func (c *Client) waitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.backend, tx)
	if err != nil {
		c.markNonceDirty()
		return nil, fmt.Errorf("wait mined tx %s: %w", tx.Hash().Hex(), err)
	}
	return receipt, nil
}
