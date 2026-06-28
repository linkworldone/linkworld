package sync

import (
	"context"
	"log"
	"math/big"
	"time"

	"linkworld-backend/internal/blockchain"
	"linkworld-backend/internal/blockchain/bindings"
	"linkworld-backend/internal/config"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// CONFIRMATIONS 是资金事件落 confirmed 前需等待的区块确认数（design §6.3/§7.0③，B5）。
//
// !!! PLACEHOLDER：K=5 为占位值，上线前须由安全/运营按 Arbitrum Sepolia 最终性确认 !!!
// 资金事件（DepositMade/DepositWithdrawn/BillPaid/BillCreated）解码后先落 seen，待块深度 ≥ K
// 才置 confirmed（IsPaid/记账生效）；非资金事件即时 confirmed。
const CONFIRMATIONS = uint64(5)

// blockRangeLimit 单次 FilterLogs 拉取的最大区块跨度（公共 RPC 多有上限）。
const blockRangeLimit = uint64(2000)

// defaultPollInterval event_sync 轮询间隔。
const defaultPollInterval = 15 * time.Second

// logSource 抽象 event_sync 依赖的链上读能力（*ethclient.Client 满足；测试用 mock 满足）。
// 拆出接口便于离线测 reorg/分页/确认深度（design §6.3 TDD：不依赖部署 Cancun 合约）。
type logSource interface {
	bind.ContractFilterer // FilterLogs / SubscribeFilterLogs（Parse* 解码所需的 backend）
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// EventSync 真实事件同步器（design §6.3）：FilterLogs 轮询 + abigen Filterer 解码 +
// (txHash,logIndex) 去重 + reorg 检测 + K 块两阶段确认 + 资金事件单一对账回填。
type EventSync struct {
	source    logSource
	contracts map[string]common.Address
	chainID   uint64

	userRepo    *repository.UserRepository
	billRepo    *repository.BillRepository
	depositRepo *repository.DepositRepository
	simRepo     *repository.SimRepository
	syncRepo    *repository.SyncStateRepository
	eventRepo   *repository.ChainEventRepository

	// 每合约一个 abigen Filterer（仅非占位零地址合约装配）。
	paymentF      *bindings.PaymentFilterer
	depositF      *bindings.DepositFilterer
	oracleF       *bindings.OracleFilterer
	userRegistryF *bindings.UserRegistryFilterer
	trafficNFTF   *bindings.TrafficCardNFTFilterer

	// 各合约地址（占位零地址则该合约为 nil，跳过订阅）。
	addrPayment      *common.Address
	addrDeposit      *common.Address
	addrOracle       *common.Address
	addrUserRegistry *common.Address
	addrTrafficNFT   *common.Address

	confirmations uint64
	rangeLimit    uint64
	pollInterval  time.Duration
}

// NewEventSync 生产路径构造（main.go 调用，签名保持不变）：从 ethClient + userRepo + 地址表装配。
// 其余 repo 复用 userRepo 的 *gorm.DB 派生（不改 main.go）。chainID 在 Start 时读链补齐。
func NewEventSync(ethClient *ethclient.Client, userRepo *repository.UserRepository, contracts map[string]common.Address) *EventSync {
	db := userRepo.DB()
	// 自带迁移 event_sync 专属表（SyncState 游标 + ChainEvent 去重/两阶段），与功能同源、幂等，
	// 不依赖 main.go 的 AutoMigrate 清单（main.go 不在 T4 改动范围）。
	if err := db.AutoMigrate(&models.SyncState{}, &models.ChainEvent{}, &models.Sim{}); err != nil {
		log.Printf("WARN: event_sync 表迁移失败：%v", err)
	}
	return newEventSync(
		ethClient,
		userRepo,
		repository.NewBillRepository(db),
		repository.NewDepositRepository(db),
		repository.NewSimRepository(db),
		repository.NewSyncStateRepository(db),
		repository.NewChainEventRepository(db),
		contracts,
		0, // chainID 启动时读链补齐
	)
}

// newEventSync 内部构造（测试可注入 mock source + 各 repo + chainID）。
func newEventSync(
	source logSource,
	userRepo *repository.UserRepository,
	billRepo *repository.BillRepository,
	depositRepo *repository.DepositRepository,
	simRepo *repository.SimRepository,
	syncRepo *repository.SyncStateRepository,
	eventRepo *repository.ChainEventRepository,
	contracts map[string]common.Address,
	chainID uint64,
) *EventSync {
	s := &EventSync{
		source:        source,
		contracts:     contracts,
		chainID:       chainID,
		userRepo:      userRepo,
		billRepo:      billRepo,
		depositRepo:   depositRepo,
		simRepo:       simRepo,
		syncRepo:      syncRepo,
		eventRepo:     eventRepo,
		confirmations: CONFIRMATIONS,
		rangeLimit:    blockRangeLimit,
		pollInterval:  defaultPollInterval,
	}
	s.assembleFilterers()
	return s
}

// assembleFilterers 按地址表装配各合约 Filterer，占位零地址合约跳过订阅 + warn（design §6.3，用 T2 IsPlaceholder）。
func (s *EventSync) assembleFilterers() {
	mk := func(name string) *common.Address {
		addr, ok := s.contracts[name]
		if !ok || config.IsPlaceholder(addr) {
			log.Printf("WARN: event_sync 跳过合约 %s 订阅（地址未配置或占位零地址，未上链）", name)
			return nil
		}
		a := addr
		return &a
	}
	if a := mk("Payment"); a != nil {
		s.addrPayment = a
		s.paymentF, _ = bindings.NewPaymentFilterer(*a, s.source)
	}
	if a := mk("Deposit"); a != nil {
		s.addrDeposit = a
		s.depositF, _ = bindings.NewDepositFilterer(*a, s.source)
	}
	if a := mk("Oracle"); a != nil {
		s.addrOracle = a
		s.oracleF, _ = bindings.NewOracleFilterer(*a, s.source)
	}
	if a := mk("UserRegistry"); a != nil {
		s.addrUserRegistry = a
		s.userRegistryF, _ = bindings.NewUserRegistryFilterer(*a, s.source)
	}
	if a := mk("TrafficCardNFT"); a != nil {
		s.addrTrafficNFT = a
		s.trafficNFTF, _ = bindings.NewTrafficCardNFTFilterer(*a, s.source)
	}
}

// Start 启动后台轮询同步（design §6.3）。
func (s *EventSync) Start(ctx context.Context) {
	go s.run(ctx)
}

func (s *EventSync) run(ctx context.Context) {
	// 注：生产路径 chainID 传 0（main.go 未读链），SyncState 以 chainID=0 单行记录游标——
	// 当前单链部署足够；多链部署需在 main.go 注入真实 chainID（留 T5/后续）。
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	// 先跑一轮，再按 ticker 周期跑。
	if err := s.SyncOnce(ctx); err != nil {
		log.Printf("WARN: event_sync 首轮失败：%v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.SyncOnce(ctx); err != nil {
				log.Printf("WARN: event_sync 轮询失败：%v", err)
			}
		}
	}
}

// SyncOnce 执行一轮同步：reorg 检测 → 拉取新区块日志解码落库 → 提升达确认深度的 seen 事件。
// 导出供测试逐轮驱动（无 ticker 依赖）。
func (s *EventSync) SyncOnce(ctx context.Context) error {
	head, err := s.source.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	latest := head.Number.Uint64()

	// 读游标（断点续传）。
	state, err := s.syncRepo.Get(s.chainID)
	if err != nil {
		return err
	}

	from := uint64(0)
	if state != nil {
		// reorg 检测：上次记录的块哈希在当前链上是否仍一致（design §6.3）。
		if reorged, forkFrom := s.detectReorg(ctx, state); reorged {
			log.Printf("WARN: event_sync 检测到 reorg，从块 %d 回退重扫", forkFrom)
			// 删受影响区间未确认 seen 记录 + 回滚其关联记账。
			if rerr := s.rollbackFrom(forkFrom); rerr != nil {
				return rerr
			}
			from = forkFrom
		} else {
			from = state.LastBlock + 1
		}
	}
	if from > latest {
		// 无新块；仍尝试提升已达确认深度的 seen 事件（K 块后链未前进的边界）。
		return s.promoteConfirmed(latest)
	}

	// 分页拉取 [from, latest]。
	for start := from; start <= latest; start += s.rangeLimit {
		end := start + s.rangeLimit - 1
		if end > latest {
			end = latest
		}
		if err := s.scanRange(ctx, start, end); err != nil {
			return err
		}
		// 推进游标到 end（记 end 块哈希供下轮 reorg 检测）。
		endHeader, herr := s.source.HeaderByNumber(ctx, new(big.Int).SetUint64(end))
		if herr != nil {
			return herr
		}
		if serr := s.syncRepo.Save(s.chainID, end, endHeader.Hash().Hex()); serr != nil {
			return serr
		}
	}

	// 提升达确认深度（latest - K）的 seen 资金事件 → confirmed + 回填终态。
	return s.promoteConfirmed(latest)
}

// detectReorg 校验上次游标块哈希在当前链上是否仍一致。不一致 → 返回回退起点（向前回退到一段安全距离）。
func (s *EventSync) detectReorg(ctx context.Context, state *models.SyncState) (bool, uint64) {
	if state.BlockHash == "" {
		return false, 0
	}
	header, err := s.source.HeaderByNumber(ctx, new(big.Int).SetUint64(state.LastBlock))
	if err != nil || header == nil {
		// 该块号在当前链上已不存在（深度回退）→ 视为 reorg，从确认深度外回退重扫。
		return true, s.safeRewindFloor(state.LastBlock)
	}
	if header.Hash().Hex() != state.BlockHash {
		// 同高度不同哈希 → reorg。
		return true, s.safeRewindFloor(state.LastBlock)
	}
	return false, 0
}

// safeRewindFloor 计算回退起点：回退到 last - K（最多到 0），已确认（深度 ≥ K）的不回退。
func (s *EventSync) safeRewindFloor(last uint64) uint64 {
	if last <= s.confirmations {
		return 0
	}
	return last - s.confirmations
}

// rollbackFrom 回退重扫：删块号 ≥ from 的未确认 seen 事件 + 回滚其关联 pending 记账，
// 并重置游标到 from-1（design §6.3）。
func (s *EventSync) rollbackFrom(from uint64) error {
	deleted, err := s.eventRepo.DeleteUnconfirmedFrom(from)
	if err != nil {
		return err
	}
	// 回滚被删事件关联的 pending 资金记账（design §6.3：已 confirmed 不回退）。
	db := s.userRepo.DB()
	for _, ev := range deleted {
		log.Printf("INFO: event_sync reorg 回滚未确认事件 %s#%d (%s, block %d)", ev.TxHash, ev.LogIndex, ev.EventName, ev.BlockNumber)
		switch ev.EventName {
		case "DepositMade", "DepositWithdrawn":
			db.Where("tx_hash = ? AND status = ?", ev.TxHash, models.DepositStatusPending).
				Delete(&models.Deposit{})
		case "SimRedeemed":
			// 回滚未确认 SIM 兑换记账（已 confirmed 不回退）。
			db.Where("tx_hash = ? AND status = ?", ev.TxHash, models.SimStatusPending).
				Delete(&models.Sim{})
		}
	}
	// 游标回退到 from-1（下轮从 from 重扫）。
	prevHash := ""
	if from > 0 {
		if h, herr := s.source.HeaderByNumber(context.Background(), new(big.Int).SetUint64(from-1)); herr == nil && h != nil {
			prevHash = h.Hash().Hex()
		}
	}
	if from == 0 {
		return s.syncRepo.Save(s.chainID, 0, prevHash)
	}
	return s.syncRepo.Save(s.chainID, from-1, prevHash)
}

// scanRange 拉取 [from,to] 区间日志，按合约地址过滤后逐条 abigen Filterer 解码落库。
func (s *EventSync) scanRange(ctx context.Context, from, to uint64) error {
	addrs := s.subscribedAddresses()
	if len(addrs) == 0 {
		return nil // 无任何非占位合约，跳过。
	}
	q := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: addrs,
	}
	logs, err := s.source.FilterLogs(ctx, q)
	if err != nil {
		return err
	}
	for _, lg := range logs {
		if lg.Removed {
			continue // reorg 移除的日志由 reorg 检测路径处理。
		}
		if err := s.dispatch(lg); err != nil {
			log.Printf("WARN: event_sync 处理日志失败 %s#%d: %v", lg.TxHash.Hex(), lg.Index, err)
		}
	}
	return nil
}

// subscribedAddresses 返回所有非占位合约地址（FilterLogs 地址过滤；占位零地址不参与）。
func (s *EventSync) subscribedAddresses() []common.Address {
	var addrs []common.Address
	for _, a := range []*common.Address{s.addrPayment, s.addrDeposit, s.addrOracle, s.addrUserRegistry, s.addrTrafficNFT} {
		if a != nil {
			addrs = append(addrs, *a)
		}
	}
	return addrs
}

// dispatch 按 (合约地址, topic0) 分发到具体 process；幂等去重 (txHash, logIndex)。
func (s *EventSync) dispatch(lg types.Log) error {
	if len(lg.Topics) == 0 {
		return nil
	}
	// (txHash, logIndex) 去重：已落库直接跳过（轮询区块重叠时防重复记账，design §6.3）。
	if s.eventRepo.Seen(lg.TxHash.Hex(), lg.Index) {
		return nil
	}
	topic0 := lg.Topics[0]

	switch {
	case s.addrUserRegistry != nil && lg.Address == *s.addrUserRegistry && topic0 == blockchain.UserRegisteredTopic:
		return s.processUserRegistered(lg)

	case s.addrDeposit != nil && lg.Address == *s.addrDeposit && topic0 == blockchain.DepositMadeTopic:
		return s.processDepositMade(lg)
	case s.addrDeposit != nil && lg.Address == *s.addrDeposit && topic0 == blockchain.DepositWithdrawnTopic:
		return s.processDepositWithdrawn(lg)
	case s.addrDeposit != nil && lg.Address == *s.addrDeposit && topic0 == blockchain.TrafficCardMintedTopic:
		return s.processTrafficCardMinted(lg)

	case s.addrPayment != nil && lg.Address == *s.addrPayment && topic0 == blockchain.BillCreatedTopic:
		return s.processBillCreated(lg)
	case s.addrPayment != nil && lg.Address == *s.addrPayment && topic0 == blockchain.BillPaidTopic:
		return s.processBillPaid(lg)
	case s.addrPayment != nil && lg.Address == *s.addrPayment && topic0 == blockchain.TrafficCardAppliedTopic:
		return s.processTrafficCardApplied(lg)

	case s.addrOracle != nil && lg.Address == *s.addrOracle && topic0 == blockchain.UsageDataSubmittedTopic:
		return s.processUsageDataSubmitted(lg)

	case s.addrTrafficNFT != nil && lg.Address == *s.addrTrafficNFT && topic0 == blockchain.CardMintedTopic:
		return s.processCardMinted(lg)
	case s.addrTrafficNFT != nil && lg.Address == *s.addrTrafficNFT && topic0 == blockchain.SimRedeemedTopic:
		return s.processSimRedeemed(lg)
	case s.addrTrafficNFT != nil && lg.Address == *s.addrTrafficNFT && topic0 == blockchain.ESimRedeemedTopic:
		return s.processESimRedeemed(lg)
	}
	return nil
}

// recordEvent 幂等落事件去重记录（status 由 fund 决定：资金事件 seen 待确认，非资金事件即 confirmed）。
// onChainRef 携带确认阶段回填所需的链上引用（如 BillPaid 的 billId），无则传 0。
func (s *EventSync) recordEvent(name string, lg types.Log, fund bool, onChainRef uint64) error {
	status := models.ChainEventStatusConfirmed
	if fund {
		status = models.ChainEventStatusSeen
	}
	return s.eventRepo.Record(&models.ChainEvent{
		TxHash:      lg.TxHash.Hex(),
		LogIndex:    lg.Index,
		EventName:   name,
		BlockNumber: lg.BlockNumber,
		BlockHash:   lg.BlockHash.Hex(),
		Status:      status,
		OnChainRef:  onChainRef,
	})
}

// --- process handlers（一律走 abigen Filterer.Parse* 解码，design §6.3）---

// processUserRegistered 解析 email/tokenId/真实时间落库（非资金事件，即时 confirmed）。
func (s *EventSync) processUserRegistered(lg types.Log) error {
	ev, err := s.userRegistryF.ParseUserRegistered(lg)
	if err != nil {
		return err
	}
	user := &models.User{
		WalletAddr:   ev.User.Hex(),
		Email:        ev.Email,
		TokenID:      uint(ev.TokenId.Uint64()),
		IsActive:     true,
		RegisteredAt: time.Now(),
	}
	if err := s.userRepo.CreateIfNotExists(user); err != nil {
		return err
	}
	return s.recordEvent("UserRegistered", lg, false, 0)
}

// processDepositMade 押金落库（资金事件，6 位精度，等 K 块确认才计入余额）。
func (s *EventSync) processDepositMade(lg types.Log) error {
	ev, err := s.depositF.ParseDepositMade(lg)
	if err != nil {
		return err
	}
	// 关联 user（链上以 wallet 为主键，DB 以 userID）。找不到 user 仍记 ChainEvent，待用户注册后由后续对账补。
	uid := s.userIDByWallet(ev.User)
	if uid != 0 {
		// 优先对账回填已有 HTTP 意向记录（无 txHash 的 pending），找不到才新建。
		// 避免「意向 + 链上事件」产生两条记录、意向那条永久卡 pending（confirmDeposit 只按 txHash 匹配）。
		s.reconcileOrCreate(uid, "deposit", ev.Amount.String(), lg)
	}
	return s.recordEvent("DepositMade", lg, true, 0)
}

// processDepositWithdrawn 提现记账（资金事件，B3 唯一 withdraw 记账路径）。
func (s *EventSync) processDepositWithdrawn(lg types.Log) error {
	ev, err := s.depositF.ParseDepositWithdrawn(lg)
	if err != nil {
		return err
	}
	uid := s.userIDByWallet(ev.User)
	if uid != 0 {
		// 提现金额 = amount（单笔 tranche 本金，合约已不再分发利息），按 6 位最小单位记 withdraw 负向记账。
		// 同样优先回填已有 withdraw 意向（无 txHash 的 pending），找不到才新建（防重复 + 永久 pending）。
		s.reconcileOrCreate(uid, "withdraw", ev.Amount.String(), lg)
	}
	return s.recordEvent("DepositWithdrawn", lg, true, 0)
}

// processBillCreated 回填链上 billId + TxHash（资金事件）。第三参 totalAmount 含费总额（design §6.3）。
func (s *EventSync) processBillCreated(lg types.Log) error {
	ev, err := s.paymentF.ParseBillCreated(lg)
	if err != nil {
		return err
	}
	uid := s.userIDByWallet(ev.User)
	if uid != 0 {
		// 回填 OnChainBillID（operatorId 不在该事件中，用 user 关联尚未关联的最新 bill）。
		// TotalAmount = amount + platformFee（含费总额），对账校验留作 processUsageDataSubmitted 维度。
		s.billRepo.SetOnChainID(ev.BillId.Uint64(), lg.TxHash.Hex(), uid, 0)
		_ = ev.TotalAmount
		_ = ev.PlatformFee
	}
	return s.recordEvent("BillCreated", lg, true, ev.BillId.Uint64())
}

// processBillPaid 唯一置 IsPaid 的路径（资金事件，B2，等 K 块确认后由 promoteConfirmed 生效）。
func (s *EventSync) processBillPaid(lg types.Log) error {
	ev, err := s.paymentF.ParseBillPaid(lg)
	if err != nil {
		return err
	}
	// 实际 IsPaid 置位推迟到 promoteConfirmed（深度 ≥ K），此处仅记 seen 事件 + 携带 billId 供确认回填。
	return s.recordEvent("BillPaid", lg, true, ev.BillId.Uint64())
}

// processUsageDataSubmitted 对账（后端计价金额 == 链上入账金额，不一致告警）。
// 只有 user indexed；operatorId/amount 在 data 区（design §6.3 解码歧义澄清）。非资金终态，记录即可。
func (s *EventSync) processUsageDataSubmitted(lg types.Log) error {
	ev, err := s.oracleF.ParseUsageDataSubmitted(lg)
	if err != nil {
		return err
	}
	uid := s.userIDByWallet(ev.User)
	if uid != 0 {
		bills, _ := s.billRepo.FindByUserID(uid)
		for _, b := range bills {
			if b.OperatorID == uint(ev.OperatorId.Uint64()) {
				if onchain := ev.Amount.String(); b.Amount != "" && b.Amount != onchain {
					log.Printf("WARN: 对账金额不一致 user=%s op=%d 后端=%s 链上=%s",
						ev.User.Hex(), ev.OperatorId.Uint64(), b.Amount, onchain)
				}
				break
			}
		}
	}
	return s.recordEvent("UsageDataSubmitted", lg, false, 0)
}

// processTrafficCardMinted Deposit 合约发卡事件落库（非资金终态，即时 confirmed）。
func (s *EventSync) processTrafficCardMinted(lg types.Log) error {
	if _, err := s.depositF.ParseTrafficCardMinted(lg); err != nil {
		return err
	}
	return s.recordEvent("TrafficCardMinted", lg, false, 0)
}

// processCardMinted TrafficCardNFT 合约发卡事件落库。
func (s *EventSync) processCardMinted(lg types.Log) error {
	if _, err := s.trafficNFTF.ParseCardMinted(lg); err != nil {
		return err
	}
	return s.recordEvent("CardMinted", lg, false, 0)
}

// processSimRedeemed 流量卡销毁兑换 SIM（新玩法，资金类/重要事件，等 K 块确认）。
// 解析 SimRedeemed(user, daysCount, tokenIds)，按 txHash 对账回填已有 pending SIM 意向：
// 命中 pending 意向（txHash 匹配 / 或无 txHash）→ 回填 txHash 保持 pending（promoteConfirmed 确认）；
// 找不到（用户直接链上 redeem、无前端意向）→ 用事件 daysCount 新建一条 pending SIM。
// 实际 confirmed 终态由 promoteConfirmed（深度 ≥ K）置位。
func (s *EventSync) processSimRedeemed(lg types.Log) error {
	ev, err := s.trafficNFTF.ParseSimRedeemed(lg)
	if err != nil {
		return err
	}
	uid := s.userIDByWallet(ev.User)
	if uid != 0 {
		s.reconcileOrCreateSim(uid, uint(ev.DaysCount.Uint64()), lg)
	}
	return s.recordEvent("SimRedeemed", lg, true, 0)
}

// processESimRedeemed eSIM 兑换事件落库，更新 SIM 记录的激活码和激活链接。
// 解析 ESimRedeemed(user, tokenId, activationCode, smDpAddress)，按 txHash 关联 SIM 记录。
func (s *EventSync) processESimRedeemed(lg types.Log) error {
	ev, err := s.trafficNFTF.ParseESimRedeemed(lg)
	if err != nil {
		return err
	}
	txHash := lg.TxHash.Hex()
	activationCode := ev.ActivationCode
	smDpAddress := ev.SmDpAddress
	activationURL := "LPKG://" + activationCode + "@" + smDpAddress
	s.simRepo.UpdateByTxHash(txHash, activationCode, activationURL, ev.TokenId.Uint64())
	return s.recordEvent("ESimRedeemed", lg, false, 0)
}

// reconcileOrCreateSim SIM 兑换事件落库的单一对账回填（避免「意向 + 链上事件」双记录）。
// 查同 user、status=pending、且 (txHash 匹配本 tx / 或 txHash 为空) 的最近一条意向：
//   - 命中 → 回填 txHash（保持 pending，交 promoteConfirmed 确认）；
//   - 未命中（链上直接 redeem 无前端意向）→ 用事件 daysCount 新建一条 pending SIM。
func (s *EventSync) reconcileOrCreateSim(uid uint, days uint, lg types.Log) {
	db := s.userRepo.DB()
	txHash := lg.TxHash.Hex()
	var intent models.Sim
	err := db.Where("user_id = ? AND status = ? AND (tx_hash = ? OR tx_hash = '' OR tx_hash IS NULL)",
		uid, models.SimStatusPending, txHash).
		Order("created_at DESC").
		First(&intent).Error
	if err == nil {
		db.Model(&models.Sim{}).Where("id = ?", intent.ID).
			Update("tx_hash", txHash)
		return
	}
	// 未命中意向 → 用事件 daysCount 新建一条带 txHash 的 pending SIM。
	s.simRepo.Create(&models.Sim{
		UserID: uid,
		Days:   days,
		TxHash: txHash,
		Status: models.SimStatusPending,
	})
}

// processTrafficCardApplied 桩事件，仅记录不改金额（design §4.3/§5.2）。
func (s *EventSync) processTrafficCardApplied(lg types.Log) error {
	if _, err := s.paymentF.ParseTrafficCardApplied(lg); err != nil {
		return err
	}
	return s.recordEvent("TrafficCardApplied", lg, false, 0)
}

// promoteConfirmed 将块深度 ≥ K 的 seen 资金事件置 confirmed，并据事件类型回填资金终态
// （IsPaid / deposit-withdraw confirmed）。design §6.3/§7.0③ B5 两阶段确认收口。
func (s *EventSync) promoteConfirmed(latest uint64) error {
	if latest < s.confirmations {
		return nil
	}
	confirmedUpTo := latest - s.confirmations
	seen, err := s.eventRepo.FindSeenUpTo(confirmedUpTo)
	if err != nil {
		return err
	}
	for _, ev := range seen {
		switch ev.EventName {
		case "BillPaid":
			// 唯一置 IsPaid 的路径（B2）：用 seen 阶段携带的链上 billId 回填（深度 ≥ K 才生效）。
			s.billRepo.MarkPaidByOnChainID(ev.OnChainRef, ev.TxHash)
		case "DepositMade":
			s.confirmDeposit(ev, "deposit")
		case "DepositWithdrawn":
			s.confirmDeposit(ev, "withdraw")
		case "BillCreated":
			// OnChainBillID 已在 seen 阶段回填，确认即终态，无额外动作。
		case "SimRedeemed":
			// SIM 兑换确认：把该 tx 对应的 pending SIM 置 confirmed（K 块后生效）。
			s.confirmSim(ev.TxHash)
		}
		if cerr := s.eventRepo.MarkConfirmed(ev.TxHash, ev.LogIndex); cerr != nil {
			return cerr
		}
	}
	return nil
}

// reconcileOrCreate 资金事件落库的单一对账回填（充值/提现重复 bug 修复）。
//
// 来源二选一防双记录：
//  1. HTTP 意向（depositService.Deposit）已先建一条 {type, amount, status:pending, tx_hash:""}（无 txHash）；
//  2. 本链上事件想再落一条带 txHash 的 pending。
//
// 若无脑 Create 会得到两条 pending；而 confirmDeposit 只按 tx_hash 匹配 → 意向那条永久卡 pending。
// 故此处优先查找同 user、同 type、同 amount、status=pending、且 tx_hash 为空的最近一条意向记录：
//   - 命中 → UPDATE 回填 tx_hash + block_hash（仍 pending，后续 promoteConfirmed 按 txHash 确认它）；
//   - 未命中（用户直接链上充值、无前端意向）→ 才 Create 新记录。
//
// 一笔充值 = 一条记录。amount 按字符串等值匹配（事件 .String() 与意向存储格式一致，均 6 位最小单位）。
func (s *EventSync) reconcileOrCreate(uid uint, typ, amount string, lg types.Log) {
	db := s.userRepo.DB()
	var intent models.Deposit
	err := db.Where("user_id = ? AND type = ? AND amount = ? AND status = ? AND (tx_hash = '' OR tx_hash IS NULL)",
		uid, typ, amount, models.DepositStatusPending).
		Order("created_at DESC").
		First(&intent).Error
	if err == nil {
		// 命中意向 → 回填 txHash/blockHash（状态保持 pending，交给 promoteConfirmed 确认）。
		db.Model(&models.Deposit{}).Where("id = ?", intent.ID).
			Updates(map[string]interface{}{
				"tx_hash":    lg.TxHash.Hex(),
				"block_hash": lg.BlockHash.Hex(),
			})
		return
	}
	// 未命中意向（含 ErrRecordNotFound 与其他查询错误）→ 新建一条带 txHash 的 pending 记录。
	s.depositRepo.Create(&models.Deposit{
		UserID:    uid,
		Amount:    amount, // 6 位最小单位
		Type:      typ,
		TxHash:    lg.TxHash.Hex(),
		Status:    models.DepositStatusPending,
		BlockHash: lg.BlockHash.Hex(),
	})
}

// confirmDeposit 将该 tx 对应的 pending deposit/withdraw 记录置 confirmed（计入余额）。
func (s *EventSync) confirmDeposit(ev models.ChainEvent, typ string) {
	db := s.userRepo.DB()
	db.Model(&models.Deposit{}).
		Where("tx_hash = ? AND type = ? AND status = ?", ev.TxHash, typ, models.DepositStatusPending).
		Update("status", models.DepositStatusConfirmed)
}

// confirmSim 将该 tx 对应的 pending SIM 记录置 confirmed（同押金两阶段，深度 ≥ K 才生效）。
func (s *EventSync) confirmSim(txHash string) {
	db := s.userRepo.DB()
	db.Model(&models.Sim{}).
		Where("tx_hash = ? AND status = ?", txHash, models.SimStatusPending).
		Update("status", models.SimStatusConfirmed)
}

// userIDByWallet 由钱包地址查 userID；找不到返回 0（调用方据此决定是否落业务行）。
func (s *EventSync) userIDByWallet(wallet common.Address) uint {
	u, err := s.userRepo.FindByWallet(wallet.Hex())
	if err != nil || u == nil {
		return 0
	}
	return u.ID
}
