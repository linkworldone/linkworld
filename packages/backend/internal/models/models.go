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

// Sim 是「流量卡销毁兑换 SIM」的链下记录（新玩法）：用户一次性销毁多张流量卡 NFT
// （每张=1 天）兑换到一张无限流量 SIM（链下实体/eSIM，配送）。
// 链上只销毁卡 + emit SimRedeemed(user, daysCount, tokenIds)；SIM 身份/天数/收件信息全记后端 DB。
// Status 两阶段（同押金）：pending(HTTP claim 意向) → confirmed(链上 SimRedeemed 等 K 块确认)。
type Sim struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"index" json:"userId"`
	TokenId        uint      `json:"tokenId"`                   // NFT tokenId (for eSIM lookup)
	Days           uint      `json:"days"`                      // 无限流量天数(=销毁卡数)
	Destination    string    `gorm:"size:8" json:"destination"` // 目的地 code(US/JP/...)
	Recipient      string    `json:"recipient"`                 // 收件人（仅 physical 邮寄需要；esim 可空）
	AddressLine    string    `json:"addressLine"`               // 收件地址（仅 physical 邮寄需要；esim 可空）
	TxHash         string    `json:"txHash"`
	Status         string    `gorm:"size:16;index" json:"status"` // pending → confirmed(同押金两阶段)
	DeliveryType   string    `gorm:"size:8;default:physical" json:"deliveryType"`
	ActivationCode string    `json:"activationCode"`              // eSIM 激活码（来自链上 ESimRedeemed 事件）
	ActivationURL  string    `json:"activationUrl"`               // eSIM 激活链接（链上 SM-DP + activationCode 组合）
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Sim.Status 取值常量（同押金两阶段状态机）。
const (
	SimStatusPending   = "pending"
	SimStatusConfirmed = "confirmed"
)

// Sim.DeliveryType 取值常量。
const (
	SimDeliveryEsim     = "esim"
	SimDeliveryPhysical = "physical"
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

// SettlementBatch 是月度结算的分批幂等记录（design §4.1/§6.1，arch-review B1：幂等键 month+batchIndex）。
// 月度结算把 amounts[] 切片为每批 ≤25 user 逐批发交易；每批以 (Month, BatchIndex) 唯一标识：
//   - confirmed 批不重发（重复触发幂等）；
//   - failed 批可重试（一批失败不影响其他批，失败批续跑）；
//   - pending_review 批因 L2 熔断/绝对闸阻断，待人工放行后再发（不自动发）。
//
// TotalAmount 落该批 sum(amounts)（6 位最小单位字符串），供历史月均熔断计算。
type SettlementBatch struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Month       string `gorm:"size:7;uniqueIndex:idx_month_batch" json:"month"` // "YYYY-MM"
	BatchIndex  int    `gorm:"uniqueIndex:idx_month_batch" json:"batch_index"`
	UserCount   int    `json:"user_count"`
	TotalAmount string `json:"total_amount"` // 6 位最小单位字符串（design §4.4）
	// Status：pending(已建未发) → confirmed(receipt.Status=1) / failed(receipt.Status=0 或发送失败，可重试)
	//        / pending_review(L2 熔断或绝对闸阻断，待人工放行)。
	Status    string    `gorm:"size:20;index" json:"status"`
	TxHash    string    `json:"tx_hash,omitempty"`
	Note      string    `json:"note,omitempty"` // 阻断原因等
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettlementBatch.Status 取值常量。
const (
	BatchStatusPending       = "pending"
	BatchStatusConfirmed     = "confirmed"
	BatchStatusFailed        = "failed"
	BatchStatusPendingReview = "pending_review"
)

// WalletNonce 是 WalletAuth 的服务端一次性 nonce 台账（arch-review 🔴 N1 红线）：
// 防签名重放——每个 nonce 由 GET /api/auth/nonce/:wallet 签发（Used=false 落库），
// WalletAuth 校验通过后置 Used=true（消费式，签过即作废）。同 nonce 二次提交必拒（WALLET-02）。
// 绑定 Wallet（小写归一化）避免跨钱包借用 nonce；ExpiresAt 兜底过期清理（窗口仅作 DoS 防护，
// 不作为防重放唯一手段——防重放靠 Used 标记，禁纯 timestamp 时间窗）。
type WalletNonce struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Wallet    string    `gorm:"index;size:42" json:"wallet"`
	Nonce     string    `gorm:"uniqueIndex;size:80" json:"nonce"`
	Used      bool      `gorm:"default:false;index" json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailVerification 是「邮箱绑定」的一次性验证码台账（与 WalletNonce 同安全模式）：
// 用户在 settings 提交邮箱 → Issue 签发 6 位数字码（crypto/rand）落库（Used=false）→
// 发码邮件 → 用户输入码 Verify 通过后置 Used=true（消费式，一次性），再由上层更新 User.Email。
// Wallet 小写归一化（绑定签发者，防跨钱包借码）；ExpiresAt 10 分钟过期；
// Attempts 记录该记录的错误尝试次数（≥5 锁死，防爆破）。
type EmailVerification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Wallet    string    `gorm:"index;size:42" json:"wallet"`
	Email     string    `gorm:"index" json:"email"`
	Code      string    `gorm:"size:6" json:"-"` // 6 位数字码；不出 JSON（防泄露）
	Used      bool      `gorm:"default:false;index" json:"used"`
	Attempts  int       `gorm:"default:0" json:"attempts"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
