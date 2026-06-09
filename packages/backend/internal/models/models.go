package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	WalletAddr   string    `gorm:"uniqueIndex;size:42" json:"wallet_addr"`
	Email        string    `gorm:"index" json:"email"`
	TokenID      uint      `json:"token_id"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	RegisteredAt time.Time `json:"registered_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Operator struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	Name            string    `json:"name"`
	Region          string    `json:"region"`
	CountryCode     string    `gorm:"size:3" json:"country_code"`
	RequiredDeposit string    `json:"required_deposit"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Deposit struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	UserID uint   `gorm:"index" json:"user_id"`
	Amount string `json:"amount"`                              // 6 位最小单位字符串（design §4.4）
	Type   string `gorm:"size:16;default:deposit" json:"type"` // "deposit" or "withdraw"
	TxHash string `json:"tx_hash,omitempty"`
	// Status 两阶段状态机（design §4.3/§4.4，B3/B5）：pending(HTTP 意向) → confirmed(链上事件等 K 块)。
	// 仅 confirmed 计入余额（GetTotalByUserID）。空字符串视为历史 confirmed（向后兼容旧数据）。
	Status string `gorm:"size:16;index" json:"status"`
	// BlockHash 该事件所在区块哈希，供 reorg 检测/回退重扫（design §6.3）。
	BlockHash string    `gorm:"size:66" json:"block_hash,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Deposit.Status 取值常量。
const (
	DepositStatusPending   = "pending"
	DepositStatusConfirmed = "confirmed"
)

type Bill struct {
	ID                   uint   `gorm:"primarykey" json:"id"`
	UserID               uint   `gorm:"index" json:"user_id"`
	OperatorID           uint   `gorm:"index" json:"operator_id"`
	Amount               string `json:"amount"`
	PlatformFee          string `json:"platform_fee"`
	TrafficCardDeduction string `json:"traffic_card_deduction"` // 本轮恒 "0"（design §4.3，合约 applyTrafficCardToBill 桩）
	// OnChainBillID 链上账单 ID，由 event_sync 监听 BillCreated 回填（design §4.4）。
	OnChainBillID uint64 `gorm:"index" json:"on_chain_bill_id"`
	// PayIntentTxHash 用户「我已发起支付」意向（与 IsPaid 解耦，design §4.3/B2）；
	// IsPaid 终态唯一由 event_sync 监听 BillPaid 置，HTTP 端点绝不置 IsPaid。
	PayIntentTxHash string     `json:"pay_intent_tx_hash,omitempty"`
	IsPaid          bool       `gorm:"default:false" json:"is_paid"`
	TxHash          string     `json:"tx_hash,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
}

type UserService struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex" json:"user_id"`
	OperatorID    uint      `gorm:"index" json:"operator_id"`
	VirtualNumber string    `json:"virtual_number"`
	Password      string    `json:"password"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	ActivatedAt   time.Time `json:"activated_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UsageData struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	OperatorID uint      `gorm:"index" json:"operator_id"`
	DataUsage  uint64    `json:"data_usage"` // 单位 MB（design §4.4）
	CallUsage  uint64    `json:"call_usage"` // 单位 minute（design §4.4）
	Timestamp  time.Time `json:"timestamp"`
	Signature  string    `gorm:"type:text" json:"signature"`
	CreatedAt  time.Time `json:"created_at"`
}

// SyncState 是 event_sync 的断点续传游标（design §6.3）：记录已扫描到的最新区块号 + 区块哈希，
// 供重启续传与 reorg 父哈希连续性检测。单 chainID 一行（ChainID 唯一）。
type SyncState struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ChainID   uint64    `gorm:"uniqueIndex" json:"chain_id"`
	LastBlock uint64    `json:"last_block"`
	BlockHash string    `gorm:"size:66" json:"block_hash"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChainEvent 是链上事件的幂等去重 + 两阶段状态记录（design §6.3）：
// (TxHash, LogIndex) 唯一键防重复落库（轮询区块重叠时）；Status 走 seen → confirmed
// （资金事件等 K 块确认才置 confirmed，reorg 时删未 confirmed 的 seen 记录重扫）。
type ChainEvent struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	TxHash      string `gorm:"size:66;uniqueIndex:idx_tx_log" json:"tx_hash"`
	LogIndex    uint   `gorm:"uniqueIndex:idx_tx_log" json:"log_index"`
	EventName   string `gorm:"size:48" json:"event_name"`
	BlockNumber uint64 `gorm:"index" json:"block_number"`
	BlockHash   string `gorm:"size:66" json:"block_hash"`
	// Status：seen(已解码落库，未达确认深度) → confirmed(深度 ≥ K，资金终态已生效)。
	// 非资金事件落库即 confirmed（无资损面，design §6.3）。
	Status string `gorm:"size:16;index" json:"status"`
	// OnChainRef 携带确认阶段回填所需的链上引用（如 BillPaid 的 billId），避免确认时重取 log。
	OnChainRef uint64    `json:"on_chain_ref"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ChainEvent.Status 取值常量。
const (
	ChainEventStatusSeen      = "seen"
	ChainEventStatusConfirmed = "confirmed"
)
