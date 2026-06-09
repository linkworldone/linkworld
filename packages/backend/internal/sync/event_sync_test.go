package sync

// T4 event_sync 真实同步的 TDD 测试（先红后绿）。
//
// 链交互取舍（与 T3 一致，见 client_test.go）：geth v1.13.5 simulated.Backend 受 London 上限，
// 跑不了业务合约的 Cancun 字节码。故事件解码用「手工构造 types.Log（topics+data 按真实事件 ABI
// 编码）→ abigen Filterer.Parse* 解码 → 断言字段」；reorg/去重/两阶段状态机用 mock logSource +
// 内存 sqlite（纯 Go，无 CGO）构造 log 序列驱动 SyncOnce 测，不依赖部署 Cancun 合约。

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"linkworld-backend/internal/blockchain"
	"linkworld-backend/internal/blockchain/bindings"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- 测试基础设施 ---

var (
	testPayment      = common.HexToAddress("0x1111111111111111111111111111111111111111")
	testDeposit      = common.HexToAddress("0x2222222222222222222222222222222222222222")
	testOracle       = common.HexToAddress("0x3333333333333333333333333333333333333333")
	testUserRegistry = common.HexToAddress("0x4444444444444444444444444444444444444444")
	testTrafficNFT   = common.HexToAddress("0x5555555555555555555555555555555555555555")
	testWallet       = common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
)

func testContracts() map[string]common.Address {
	return map[string]common.Address{
		"Payment":        testPayment,
		"Deposit":        testDeposit,
		"Oracle":         testOracle,
		"UserRegistry":   testUserRegistry,
		"TrafficCardNFT": testTrafficNFT,
	}
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// 每个测试独立内存库（DSN 用测试名做命名内存库，cache=shared 但库名唯一 → 互不干扰）。
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Bill{}, &models.Deposit{},
		&models.SyncState{}, &models.ChainEvent{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// fakeSource 是 mock logSource：可编程 head/区块哈希链 + 按区块号返回的 logs。
type fakeSource struct {
	head      uint64
	hashes    map[uint64]common.Hash // 块号 → 块哈希（reorg 时改写）
	logsByBlk map[uint64][]types.Log
}

func newFakeSource() *fakeSource {
	return &fakeSource{hashes: map[uint64]common.Hash{}, logsByBlk: map[uint64][]types.Log{}}
}

func (f *fakeSource) hashFor(n uint64) common.Hash {
	if h, ok := f.hashes[n]; ok {
		return h
	}
	// 默认确定性哈希。
	return common.BigToHash(new(big.Int).SetUint64(0x1000 + n))
}

func (f *fakeSource) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var out []types.Log
	from := q.FromBlock.Uint64()
	to := q.ToBlock.Uint64()
	for b := from; b <= to; b++ {
		for _, lg := range f.logsByBlk[b] {
			lg.BlockHash = f.hashFor(b)
			out = append(out, lg)
		}
	}
	return out, nil
}

func (f *fakeSource) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, nil
}

func (f *fakeSource) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	n := f.head
	if number != nil {
		n = number.Uint64()
	}
	if n > f.head {
		return nil, ethereum.NotFound
	}
	return &types.Header{Number: new(big.Int).SetUint64(n), Extra: f.hashFor(n).Bytes()}, nil
}

// headerHashOverride：types.Header.Hash() 由内容算出，无法直接塞我们的哈希。
// 这里改用 Extra 携带哈希字节，并让 detectReorg 通过 header.Hash() 比对——
// 为使比对稳定，fakeSource 用 Extra 编码块哈希，Header.Hash() 因 Extra 不同而不同，
// 仍能区分「同高度不同内容」。测试断言用 SyncState.BlockHash 的变化来验证 reorg 路径。

func newSyncWithFake(t *testing.T, src *fakeSource, db *gorm.DB) (*EventSync, *repository.UserRepository) {
	t.Helper()
	userRepo := repository.NewUserRepository(db)
	s := newEventSync(
		src,
		userRepo,
		repository.NewBillRepository(db),
		repository.NewDepositRepository(db),
		repository.NewSyncStateRepository(db),
		repository.NewChainEventRepository(db),
		testContracts(),
		31337,
	)
	// 测试用低确认数便于驱动两阶段。
	s.confirmations = 3
	return s, userRepo
}

// --- 事件日志构造 helper（用 abigen 绑定的 ABI 编码 data + topics）---

func abiFromJSON(s string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(s))
}

// buildLog 用合约 ABI 构造一条 types.Log：indexed 入 topics，非 indexed 入 data。
func buildLog(t *testing.T, abiJSON, event string, contract common.Address, blk uint64, logIdx uint, indexed []common.Hash, nonIndexed ...interface{}) types.Log {
	t.Helper()
	parsed, err := abiFromJSON(abiJSON)
	if err != nil {
		t.Fatalf("parse abi: %v", err)
	}
	ev, ok := parsed.Events[event]
	if !ok {
		t.Fatalf("event %s not in abi", event)
	}
	data, err := ev.Inputs.NonIndexed().Pack(nonIndexed...)
	if err != nil {
		t.Fatalf("pack %s: %v", event, err)
	}
	topics := append([]common.Hash{ev.ID}, indexed...)
	return types.Log{
		Address:     contract,
		Topics:      topics,
		Data:        data,
		BlockNumber: blk,
		TxHash:      common.BigToHash(big.NewInt(int64(blk*1000 + uint64(logIdx)))),
		Index:       logIdx,
	}
}

func TestParseBillCreated_TotalAmountIncludesFee(t *testing.T) {
	// design §6.3：BillCreated 第三参 totalAmount = amount + platformFee（含费总额），勿当裸 amount。
	billID := big.NewInt(42)
	total := big.NewInt(1_025_000) // 1.025 USDT（6 位）= amount 1.0 + fee 0.025
	fee := big.NewInt(25_000)
	lg := buildLog(t, bindings.PaymentMetaData.ABI, "BillCreated", testPayment, 10, 0,
		[]common.Hash{common.BigToHash(billID), common.BytesToHash(testWallet.Bytes())},
		total, fee)

	f, err := bindings.NewPaymentFilterer(testPayment, nil)
	if err != nil {
		t.Fatal(err)
	}
	ev, err := f.ParseBillCreated(lg)
	if err != nil {
		t.Fatalf("ParseBillCreated: %v", err)
	}
	if ev.BillId.Cmp(billID) != 0 {
		t.Errorf("billId = %s, want %s", ev.BillId, billID)
	}
	if ev.User != testWallet {
		t.Errorf("user = %s, want %s", ev.User, testWallet)
	}
	if ev.TotalAmount.Cmp(total) != 0 {
		t.Errorf("totalAmount = %s, want %s (含费总额)", ev.TotalAmount, total)
	}
	if ev.PlatformFee.Cmp(fee) != 0 {
		t.Errorf("platformFee = %s, want %s", ev.PlatformFee, fee)
	}
}

func TestParseUsageDataSubmitted_OnlyUserIndexed(t *testing.T) {
	// design §6.3：UsageDataSubmitted 只有 user indexed；operatorId/amount 在 data 区。
	opID := big.NewInt(7)
	amount := big.NewInt(500_000)
	lg := buildLog(t, bindings.OracleMetaData.ABI, "UsageDataSubmitted", testOracle, 11, 1,
		[]common.Hash{common.BytesToHash(testWallet.Bytes())}, // 仅 user indexed
		opID, amount)

	f, _ := bindings.NewOracleFilterer(testOracle, nil)
	ev, err := f.ParseUsageDataSubmitted(lg)
	if err != nil {
		t.Fatalf("ParseUsageDataSubmitted: %v", err)
	}
	if ev.User != testWallet {
		t.Errorf("user = %s, want %s", ev.User, testWallet)
	}
	if ev.OperatorId.Cmp(opID) != 0 {
		t.Errorf("operatorId = %s, want %s (data 区)", ev.OperatorId, opID)
	}
	if ev.Amount.Cmp(amount) != 0 {
		t.Errorf("amount = %s, want %s (data 区，6 位精度)", ev.Amount, amount)
	}
}

func TestParseDepositWithdrawn_PrincipalPlusInterest(t *testing.T) {
	lg := buildLog(t, bindings.DepositMetaData.ABI, "DepositWithdrawn", testDeposit, 12, 0,
		[]common.Hash{common.BytesToHash(testWallet.Bytes())},
		big.NewInt(1_000_000), big.NewInt(3_000)) // principal 1.0 + interest 0.003
	f, _ := bindings.NewDepositFilterer(testDeposit, nil)
	ev, err := f.ParseDepositWithdrawn(lg)
	if err != nil {
		t.Fatalf("ParseDepositWithdrawn: %v", err)
	}
	if ev.Principal.Cmp(big.NewInt(1_000_000)) != 0 || ev.Interest.Cmp(big.NewInt(3_000)) != 0 {
		t.Errorf("principal/interest = %s/%s", ev.Principal, ev.Interest)
	}
}

func TestSignatureTopicsMatchBindings(t *testing.T) {
	// signatures.go 的 topic0 必须与 abigen 绑定 event ID 一致（防 BillCreated 再次写错参数个数）。
	cases := []struct {
		name     string
		abiJSON  string
		event    string
		expected common.Hash
	}{
		{"BillCreated", bindings.PaymentMetaData.ABI, "BillCreated", blockchain.BillCreatedTopic},
		{"BillPaid", bindings.PaymentMetaData.ABI, "BillPaid", blockchain.BillPaidTopic},
		{"TrafficCardApplied", bindings.PaymentMetaData.ABI, "TrafficCardApplied", blockchain.TrafficCardAppliedTopic},
		{"DepositMade", bindings.DepositMetaData.ABI, "DepositMade", blockchain.DepositMadeTopic},
		{"DepositWithdrawn", bindings.DepositMetaData.ABI, "DepositWithdrawn", blockchain.DepositWithdrawnTopic},
		{"TrafficCardMinted", bindings.DepositMetaData.ABI, "TrafficCardMinted", blockchain.TrafficCardMintedTopic},
		{"UsageDataSubmitted", bindings.OracleMetaData.ABI, "UsageDataSubmitted", blockchain.UsageDataSubmittedTopic},
		{"UserRegistered", bindings.UserRegistryMetaData.ABI, "UserRegistered", blockchain.UserRegisteredTopic},
		{"CardMinted", bindings.TrafficCardNFTMetaData.ABI, "CardMinted", blockchain.CardMintedTopic},
		{"ServiceActivated", bindings.TrafficCardNFTMetaData.ABI, "ServiceActivated", blockchain.ServiceActivatedTopic},
	}
	for _, c := range cases {
		parsed, err := abiFromJSON(c.abiJSON)
		if err != nil {
			t.Fatalf("%s parse abi: %v", c.name, err)
		}
		if parsed.Events[c.event].ID != c.expected {
			t.Errorf("%s topic0 = %s, signatures.go = %s", c.name, parsed.Events[c.event].ID, c.expected)
		}
	}
}

// --- 落库 / 去重 / 两阶段 / reorg（mock source + sqlite）---

func seedUser(t *testing.T, ur *repository.UserRepository) uint {
	t.Helper()
	u := &models.User{WalletAddr: testWallet.Hex(), IsActive: true}
	if err := ur.Create(u); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u.ID
}

func TestSyncOnce_DepositMade_TwoPhaseConfirmation(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	// 块 10 一条 DepositMade（资金事件）。head=10 → 深度不足 K(3)，应为 pending。
	src.head = 10
	src.logsByBlk[10] = []types.Log{
		buildLog(t, bindings.DepositMetaData.ABI, "DepositMade", testDeposit, 10, 0,
			[]common.Hash{common.BytesToHash(testWallet.Bytes())}, big.NewInt(1_000_000)),
	}
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	depRepo := repository.NewDepositRepository(db)
	deps, _ := depRepo.FindByUserID(uid)
	if len(deps) != 1 || deps[0].Status != models.DepositStatusPending {
		t.Fatalf("期望 1 条 pending deposit，得到 %+v", deps)
	}
	// pending 不计入余额。
	if total, _ := depRepo.GetTotalByUserID(uid); total != "0" && total != "" {
		t.Errorf("pending 不应计入余额，total=%s", total)
	}

	// 链推进到 13（深度 = 13-10 = 3 ≥ K）→ 确认。
	src.head = 13
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	deps, _ = depRepo.FindByUserID(uid)
	if deps[0].Status != models.DepositStatusConfirmed {
		t.Fatalf("深度达 K 后应 confirmed，得到 %s", deps[0].Status)
	}
	if total, _ := depRepo.GetTotalByUserID(uid); total != "1000000" {
		t.Errorf("confirmed 后应计入余额 1000000，得到 %s", total)
	}
}

func TestSyncOnce_Dedup_NoDoubleLedger(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	src.head = 5
	src.logsByBlk[5] = []types.Log{
		buildLog(t, bindings.DepositMetaData.ABI, "DepositMade", testDeposit, 5, 0,
			[]common.Hash{common.BytesToHash(testWallet.Bytes())}, big.NewInt(2_000_000)),
	}
	// 第一轮落库。
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	// 强制重扫同一区块（模拟轮询区块重叠）：重置游标。
	repository.NewSyncStateRepository(db).Save(31337, 4, "")
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	var depCount int64
	db.Model(&models.Deposit{}).Where("user_id = ?", uid).Count(&depCount)
	if depCount != 1 {
		t.Fatalf("(txHash,logIndex) 去重失败，deposit 行数 = %d，应为 1", depCount)
	}
	var evCount int64
	db.Model(&models.ChainEvent{}).Count(&evCount)
	if evCount != 1 {
		t.Fatalf("ChainEvent 去重失败，行数 = %d，应为 1", evCount)
	}
}

func TestSyncOnce_BillPaid_OnlyEventSetsIsPaid(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	// 预置一笔 DB bill（计价已出账，未付）。
	billRepo := repository.NewBillRepository(db)
	billRepo.Create(&models.Bill{UserID: uid, OperatorID: 7, Amount: "1000000", IsPaid: false})

	// 块 10：BillCreated 回填 OnChainBillID=42。
	src.head = 10
	src.logsByBlk[10] = []types.Log{
		buildLog(t, bindings.PaymentMetaData.ABI, "BillCreated", testPayment, 10, 0,
			[]common.Hash{common.BigToHash(big.NewInt(42)), common.BytesToHash(testWallet.Bytes())},
			big.NewInt(1_000_000), big.NewInt(25_000)),
	}
	// 块 11：BillPaid（资金事件）。
	src.logsByBlk[11] = []types.Log{
		buildLog(t, bindings.PaymentMetaData.ABI, "BillPaid", testPayment, 11, 0,
			[]common.Hash{common.BigToHash(big.NewInt(42)), common.BytesToHash(testWallet.Bytes())},
			big.NewInt(1_025_000), big.NewInt(1_000_000)),
	}
	src.head = 11 // BillPaid 深度不足 → 暂不置 IsPaid。
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	var b models.Bill
	db.First(&b, "user_id = ?", uid)
	if b.OnChainBillID != 42 {
		t.Errorf("BillCreated 应回填 OnChainBillID=42，得到 %d", b.OnChainBillID)
	}
	if b.IsPaid {
		t.Errorf("BillPaid 深度不足 K 时不应置 IsPaid")
	}

	// 推进到 14（BillPaid 深度 = 14-11 = 3 ≥ K）→ IsPaid 生效。
	src.head = 14
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	db.First(&b, "user_id = ?", uid)
	if !b.IsPaid {
		t.Errorf("BillPaid 确认后应置 IsPaid=true")
	}
	if b.PaidAt == nil {
		t.Errorf("IsPaid 时应有 PaidAt")
	}
}

func TestSyncOnce_Reorg_RollbackUnconfirmedSeen(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	// 块 10 一条 DepositMade（seen，未确认）；head=11，深度不足 K(3)。
	src.head = 11
	src.logsByBlk[10] = []types.Log{
		buildLog(t, bindings.DepositMetaData.ABI, "DepositMade", testDeposit, 10, 0,
			[]common.Hash{common.BytesToHash(testWallet.Bytes())}, big.NewInt(9_000_000)),
	}
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	var seenCount int64
	db.Model(&models.ChainEvent{}).Where("status = ?", models.ChainEventStatusSeen).Count(&seenCount)
	if seenCount != 1 {
		t.Fatalf("首轮应有 1 条 seen 事件，得到 %d", seenCount)
	}

	// 模拟 reorg：分叉点在块 10，块 10 及之后（含游标块 11）哈希全部改变，
	// 且新链上块 10 不再有那条 deposit log（被重组掉）。
	src.hashes[10] = common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000000000")
	src.hashes[11] = common.HexToHash("0xfeedface00000000000000000000000000000000000000000000000000000000")
	delete(src.logsByBlk, 10) // 重组后该 log 消失
	src.head = 12
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	// reorg 回退应删除未确认 seen 事件及其 pending 记账。
	var seenAfter, depAfter int64
	db.Model(&models.ChainEvent{}).Where("status = ?", models.ChainEventStatusSeen).Count(&seenAfter)
	db.Model(&models.Deposit{}).Where("user_id = ?", uid).Count(&depAfter)
	if seenAfter != 0 {
		t.Errorf("reorg 后未确认 seen 事件应被删除，剩 %d", seenAfter)
	}
	if depAfter != 0 {
		t.Errorf("reorg 后被重组掉的 deposit 不应残留，剩 %d", depAfter)
	}
}

func TestSyncOnce_SkipsPlaceholderContracts(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	userRepo := repository.NewUserRepository(db)
	contracts := testContracts()
	contracts["Payment"] = common.Address{} // 占位零地址
	s := newEventSync(src, userRepo, repository.NewBillRepository(db),
		repository.NewDepositRepository(db), repository.NewSyncStateRepository(db),
		repository.NewChainEventRepository(db), contracts, 31337)
	if s.addrPayment != nil {
		t.Errorf("占位零地址 Payment 应跳过订阅（addrPayment 应为 nil）")
	}
	if s.paymentF != nil {
		t.Errorf("占位零地址 Payment 不应装配 Filterer")
	}
	// 非占位合约仍装配。
	if s.addrDeposit == nil {
		t.Errorf("非占位 Deposit 应装配")
	}
}
