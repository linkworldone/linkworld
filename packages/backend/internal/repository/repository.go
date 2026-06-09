package repository

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"linkworld-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByWallet(wallet string) (*models.User, error) {
	var user models.User
	err := r.db.Where("wallet_addr = ?", wallet).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Exists(wallet string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("wallet_addr = ?", wallet).Count(&count)
	return count > 0
}

// DB 暴露底层 *gorm.DB，供 event_sync(T4) 复用同一连接构造 Bill/Deposit/Sync 等 repo，
// 不改 NewEventSync 既有签名（main.go 仅传 userRepo）。
func (r *UserRepository) DB() *gorm.DB {
	return r.db
}

// CreateIfNotExists 幂等创建用户（UserRegistered 事件可能重放）。已存在则不报错。
func (r *UserRepository) CreateIfNotExists(user *models.User) error {
	if r.Exists(user.WalletAddr) {
		return nil
	}
	return r.db.Create(user).Error
}

type OperatorRepository struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) *OperatorRepository {
	return &OperatorRepository{db: db}
}

func (r *OperatorRepository) FindAll() ([]models.Operator, error) {
	var operators []models.Operator
	err := r.db.Where("is_active = ?", true).Find(&operators).Error
	return operators, err
}

func (r *OperatorRepository) FindByID(id uint) (*models.Operator, error) {
	var operator models.Operator
	err := r.db.First(&operator, id).Error
	if err != nil {
		return nil, err
	}
	return &operator, nil
}

func (r *OperatorRepository) Create(operator *models.Operator) error {
	return r.db.Create(operator).Error
}

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(bill *models.Bill) error {
	return r.db.Create(bill).Error
}

func (r *BillRepository) FindByUserID(userID uint) ([]models.Bill, error) {
	var bills []models.Bill
	err := r.db.Where("user_id = ?", userID).Find(&bills).Error
	return bills, err
}

func (r *BillRepository) FindUnpaidByUserID(userID uint) ([]models.Bill, error) {
	var bills []models.Bill
	err := r.db.Where("user_id = ? AND is_paid = ?", userID, false).Find(&bills).Error
	return bills, err
}

func (r *BillRepository) MarkAsPaid(id uint, txHash string) error {
	now := time.Now()
	return r.db.Model(&models.Bill{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_paid": true,
		"paid_at": now,
		"tx_hash": txHash,
	}).Error
}

// SetPayIntent 写用户支付意向（B2，design §4.3/§6.6）：只更新 PayIntentTxHash，绝不动 IsPaid。
// 限定 user_id 匹配（防越权改他人账单）。找不到匹配账单返回 gorm.ErrRecordNotFound。
func (r *BillRepository) SetPayIntent(userID, billID uint, txHash string) error {
	res := r.db.Model(&models.Bill{}).
		Where("id = ? AND user_id = ?", billID, userID).
		Update("pay_intent_tx_hash", txHash)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type UserServiceRepository struct {
	db *gorm.DB
}

func NewUserServiceRepository(db *gorm.DB) *UserServiceRepository {
	return &UserServiceRepository{db: db}
}

func (r *UserServiceRepository) Create(service *models.UserService) error {
	return r.db.Create(service).Error
}

func (r *UserServiceRepository) FindByUserID(userID uint) (*models.UserService, error) {
	var service models.UserService
	err := r.db.Where("user_id = ?", userID).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *UserServiceRepository) GetActiveByUserID(userID uint) (*models.UserService, error) {
	var service models.UserService
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *UserServiceRepository) Deactivate(userID uint) error {
	return r.db.Model(&models.UserService{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

func (r *UserServiceRepository) Update(service *models.UserService) error {
	return r.db.Save(service).Error
}

type DepositRepository struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) *DepositRepository {
	return &DepositRepository{db: db}
}

func (r *DepositRepository) Create(deposit *models.Deposit) error {
	return r.db.Create(deposit).Error
}

func (r *DepositRepository) FindByUserID(userID uint) ([]models.Deposit, error) {
	var deposits []models.Deposit
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&deposits).Error
	return deposits, err
}

// GetTotalByUserID 只统计 confirmed（含历史空 status）押金（design §4.3：pending 意向不计入余额）。
func (r *DepositRepository) GetTotalByUserID(userID uint) (string, error) {
	var result struct {
		Total string
	}
	err := r.db.Model(&models.Deposit{}).
		Select("COALESCE(SUM(amount), '0') as total").
		Where("user_id = ? AND (type = ? OR type = ?) AND (status = ? OR status = ? OR status IS NULL)",
			userID, "deposit", "", models.DepositStatusConfirmed, "").
		Scan(&result).Error
	return result.Total, err
}

type UsageDataRepository struct {
	db *gorm.DB
}

func NewUsageDataRepository(db *gorm.DB) *UsageDataRepository {
	return &UsageDataRepository{db: db}
}

func (r *UsageDataRepository) Create(usage *models.UsageData) error {
	return r.db.Create(usage).Error
}

func (r *UsageDataRepository) FindByUserID(userID uint) ([]models.UsageData, error) {
	var usageData []models.UsageData
	err := r.db.Where("user_id = ?", userID).Find(&usageData).Error
	return usageData, err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// --- event_sync(T4) 回填用的扩展方法 ---

// SetOnChainID 回填链上账单 ID + TxHash（design §6.3 processBillCreated）。
// 通过 user 关联尚未关联 OnChainBillID 的最新 DB bill（BillCreated 事件不含 operatorId）；
// operatorID > 0 时叠加 operator 维度精确匹配。找不到不报错（链先于 DB 写的窗口）。
func (r *BillRepository) SetOnChainID(onChainBillID uint64, txHash string, userID, operatorID uint) error {
	q := r.db.Model(&models.Bill{}).
		Where("user_id = ? AND on_chain_bill_id = 0", userID)
	if operatorID > 0 {
		q = q.Where("operator_id = ?", operatorID)
	}
	// gorm 的 Update 不直接支持 Order+Limit 的子查询定位，先选出目标 ID 再更新（确定性取最新一条）。
	var target models.Bill
	if err := q.Order("created_at desc, id desc").First(&target).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return r.db.Model(&models.Bill{}).Where("id = ?", target.ID).
		Updates(map[string]interface{}{
			"on_chain_bill_id": onChainBillID,
			"tx_hash":          txHash,
		}).Error
}

// MarkPaidByOnChainID 唯一置 IsPaid 的路径（design §4.3/B2，仅 event_sync 确认后调用）。
func (r *BillRepository) MarkPaidByOnChainID(onChainBillID uint64, txHash string) error {
	now := time.Now()
	return r.db.Model(&models.Bill{}).
		Where("on_chain_bill_id = ?", onChainBillID).
		Updates(map[string]interface{}{
			"is_paid": true,
			"paid_at": now,
			"tx_hash": txHash,
		}).Error
}

// CreateConfirmed 幂等记一笔已确认的资金记录（deposit/withdraw，design §4.3 单一对账路径）。
// 由 (tx_hash, type, user_id) 去重，重复事件不重复落库。
func (r *DepositRepository) CreateConfirmed(d *models.Deposit) error {
	d.Status = models.DepositStatusConfirmed
	var count int64
	r.db.Model(&models.Deposit{}).
		Where("tx_hash = ? AND type = ? AND user_id = ?", d.TxHash, d.Type, d.UserID).
		Count(&count)
	if count > 0 {
		return nil
	}
	return r.db.Create(d).Error
}

// SyncStateRepository 维护 event_sync 断点续传游标（design §6.3）。
type SyncStateRepository struct {
	db *gorm.DB
}

func NewSyncStateRepository(db *gorm.DB) *SyncStateRepository {
	return &SyncStateRepository{db: db}
}

// Get 返回指定链的游标；无记录返回 (nil, nil)。
func (r *SyncStateRepository) Get(chainID uint64) (*models.SyncState, error) {
	var s models.SyncState
	err := r.db.Where("chain_id = ?", chainID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Save upsert 游标（last block + blockHash）。
func (r *SyncStateRepository) Save(chainID, lastBlock uint64, blockHash string) error {
	var s models.SyncState
	err := r.db.Where("chain_id = ?", chainID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&models.SyncState{ChainID: chainID, LastBlock: lastBlock, BlockHash: blockHash}).Error
	}
	if err != nil {
		return err
	}
	s.LastBlock = lastBlock
	s.BlockHash = blockHash
	return r.db.Save(&s).Error
}

// SettlementBatchRepository 维护月度结算分批的幂等记录（design §4.1/§6.1，幂等键 month+batchIndex）。
type SettlementBatchRepository struct {
	db *gorm.DB
}

func NewSettlementBatchRepository(db *gorm.DB) *SettlementBatchRepository {
	return &SettlementBatchRepository{db: db}
}

// Migrate 自建 settlement_batches 表（自包含，不依赖 main.go 的 AutoMigrate 清单——与 event_sync 同策略）。
func (r *SettlementBatchRepository) Migrate() error {
	return r.db.AutoMigrate(&models.SettlementBatch{})
}

// Get 返回 (month, batchIndex) 的批次记录；无记录返回 (nil, nil)。
func (r *SettlementBatchRepository) Get(month string, batchIndex int) (*models.SettlementBatch, error) {
	var b models.SettlementBatch
	err := r.db.Where("month = ? AND batch_index = ?", month, batchIndex).First(&b).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Save upsert 批次记录（按 month+batchIndex 幂等键）。
func (r *SettlementBatchRepository) Save(b *models.SettlementBatch) error {
	var existing models.SettlementBatch
	err := r.db.Where("month = ? AND batch_index = ?", b.Month, b.BatchIndex).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(b).Error
	}
	if err != nil {
		return err
	}
	b.ID = existing.ID
	b.CreatedAt = existing.CreatedAt
	return r.db.Save(b).Error
}

// HistoricalConfirmedTotals 返回 beforeMonth 之前（严格小于）已 confirmed 批次的总额清单
// （6 位最小单位字符串），供 L2 历史月均熔断计算。无样本返回空切片（冷启动回退由调用方处理）。
func (r *SettlementBatchRepository) HistoricalConfirmedTotals(beforeMonth string) ([]string, error) {
	var rows []models.SettlementBatch
	err := r.db.
		Where("month < ? AND status = ?", beforeMonth, models.BatchStatusConfirmed).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, b := range rows {
		out = append(out, b.TotalAmount)
	}
	return out, nil
}

// WalletNonceRepository 维护 WalletAuth 一次性 nonce 台账（arch-review 🔴 N1 红线，消费式防重放）。
type WalletNonceRepository struct {
	db *gorm.DB
}

func NewWalletNonceRepository(db *gorm.DB) *WalletNonceRepository {
	return &WalletNonceRepository{db: db}
}

// Migrate 自建 wallet_nonces 表（自包含，不依赖 main.go AutoMigrate 清单——与 settlement/event_sync 同策略）。
func (r *WalletNonceRepository) Migrate() error {
	return r.db.AutoMigrate(&models.WalletNonce{})
}

// nonceTTL 是 nonce 兜底过期窗口（仅 DoS 防护，防重放靠 Used 标记，非纯时间窗）。
const nonceTTL = 10 * time.Minute

// Issue 为 wallet 签发一个新的一次性 nonce（crypto/rand 32 字节十六进制），落库 Used=false。
// 返回 nonce 字符串供前端签名。wallet 统一小写归一化。
func (r *WalletNonceRepository) Issue(wallet string) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(buf)
	rec := &models.WalletNonce{
		Wallet:    strings.ToLower(wallet),
		Nonce:     nonce,
		Used:      false,
		ExpiresAt: time.Now().Add(nonceTTL),
	}
	if err := r.db.Create(rec).Error; err != nil {
		return "", err
	}
	return nonce, nil
}

// Consume 原子消费 nonce（一次性台账核心）：仅当 (nonce, wallet, used=false, 未过期) 全部满足时，
// 把 Used 置 true 并返回 true（鉴权放行）；否则返回 false（拒绝——未签发/已用/钱包不符/过期）。
// 用条件 UPDATE 的 RowsAffected 判定，保证并发下同 nonce 只被消费一次（防 WALLET-02 重放）。
func (r *WalletNonceRepository) Consume(wallet, nonce string) bool {
	res := r.db.Model(&models.WalletNonce{}).
		Where("nonce = ? AND wallet = ? AND used = ? AND expires_at > ?",
			nonce, strings.ToLower(wallet), false, time.Now()).
		Update("used", true)
	if res.Error != nil {
		return false
	}
	return res.RowsAffected == 1
}

// ChainEventRepository 维护链上事件幂等去重 + 两阶段状态（design §6.3）。
type ChainEventRepository struct {
	db *gorm.DB
}

func NewChainEventRepository(db *gorm.DB) *ChainEventRepository {
	return &ChainEventRepository{db: db}
}

// Seen 报告 (txHash, logIndex) 是否已落库（幂等去重）。
func (r *ChainEventRepository) Seen(txHash string, logIndex uint) bool {
	var count int64
	r.db.Model(&models.ChainEvent{}).Where("tx_hash = ? AND log_index = ?", txHash, logIndex).Count(&count)
	return count > 0
}

// Record 幂等落一条事件记录（已存在则不重复创建）。
func (r *ChainEventRepository) Record(e *models.ChainEvent) error {
	if r.Seen(e.TxHash, e.LogIndex) {
		return nil
	}
	return r.db.Create(e).Error
}

// PromoteConfirmed 将块号 ≤ confirmedUpTo 的 seen 事件批量置 confirmed，并返回被提升的事件
// （供调用方对每条做资金终态回填）。
func (r *ChainEventRepository) FindSeenUpTo(confirmedUpTo uint64) ([]models.ChainEvent, error) {
	var evs []models.ChainEvent
	err := r.db.Where("status = ? AND block_number <= ?", models.ChainEventStatusSeen, confirmedUpTo).
		Order("block_number asc, log_index asc").Find(&evs).Error
	return evs, err
}

// MarkConfirmed 将单条事件置 confirmed。
func (r *ChainEventRepository) MarkConfirmed(txHash string, logIndex uint) error {
	return r.db.Model(&models.ChainEvent{}).
		Where("tx_hash = ? AND log_index = ?", txHash, logIndex).
		Update("status", models.ChainEventStatusConfirmed).Error
}

// DeleteUnconfirmedFrom 回退重扫：删除块号 ≥ fromBlock 的未确认（seen）事件记录
// （design §6.3 reorg：已 confirmed 视为最终不回退）。返回被删事件供调用方回滚关联记账。
func (r *ChainEventRepository) DeleteUnconfirmedFrom(fromBlock uint64) ([]models.ChainEvent, error) {
	var evs []models.ChainEvent
	if err := r.db.Where("status = ? AND block_number >= ?", models.ChainEventStatusSeen, fromBlock).
		Find(&evs).Error; err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return nil, nil
	}
	err := r.db.Where("status = ? AND block_number >= ?", models.ChainEventStatusSeen, fromBlock).
		Delete(&models.ChainEvent{}).Error
	return evs, err
}
