package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"linkworld-backend/internal/middleware"
	"linkworld-backend/internal/repository"
	"linkworld-backend/internal/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService         *services.UserService
	operatorService     *services.OperatorService
	billingService      *services.BillingService
	oracleService       *services.OracleService
	notificationService *services.NotificationService
	depositService      *services.DepositService
	userServiceService  *services.UserServiceService
	virtualGen          *services.VirtualNumberGenerator
	oracleV2            *services.OracleServiceV2
	usageService        *services.UsageService
	simService          *services.SimService
	emailService        *services.EmailService
	// nonceRepo 供 GET /api/auth/nonce/:wallet 签发 WalletAuth 一次性 nonce（arch-review 🔴 N1）。
	nonceRepo *repository.WalletNonceRepository
}

func NewHandler(
	userService *services.UserService,
	operatorService *services.OperatorService,
	billingService *services.BillingService,
	oracleService *services.OracleService,
	notificationService *services.NotificationService,
	depositService *services.DepositService,
	userServiceService *services.UserServiceService,
	virtualGen *services.VirtualNumberGenerator,
	oracleV2 *services.OracleServiceV2,
	usageService *services.UsageService,
	simService *services.SimService,
	emailService *services.EmailService,
	nonceRepo *repository.WalletNonceRepository,
) *Handler {
	return &Handler{
		userService:         userService,
		operatorService:     operatorService,
		billingService:      billingService,
		oracleService:       oracleService,
		notificationService: notificationService,
		depositService:      depositService,
		userServiceService:  userServiceService,
		virtualGen:          virtualGen,
		oracleV2:            oracleV2,
		usageService:        usageService,
		simService:          simService,
		emailService:        emailService,
		nonceRepo:           nonceRepo,
	}
}

// GetWalletNonce 签发一个一次性 WalletAuth nonce（arch-review 🔴 N1 防重放台账）。
// 公开端点（无需鉴权）：前端取 nonce 后用钱包私钥按 EIP-712 签名 (wallet,nonce,action,chainId)，
// 再带 X-Wallet-* 头请求受 WalletAuth 保护的端点。nonce 一次性消费，签过即作废。
func (h *Handler) GetWalletNonce(c *gin.Context) {
	wallet := c.Param("wallet")
	if !common.IsHexAddress(wallet) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet address"})
		return
	}
	nonce, err := h.nonceRepo.Issue(wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nonce": nonce})
}

// BindEmailRequest 提交待绑定邮箱（settings 第一步）。wallet 不收 body，取自 WalletAuth CtxWallet。
type BindEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest 提交邮箱 + 6 位验证码完成绑定（settings 第二步）。
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// BindEmail 发起邮箱绑定：签发 6 位验证码 → 发码邮件（未配 SMTP 走 [DEV] 日志降级）。
// R1 防越权：归属 wallet 唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body。
// 限流命中（同 wallet 60 秒内重复请求）→ 429；其余错误 → 400。
func (h *Handler) BindEmail(c *gin.Context) {
	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authWallet := c.GetString(middleware.CtxWallet)
	if err := h.emailService.RequestBind(authWallet, req.Email); err != nil {
		if errors.Is(err, repository.ErrEmailRateLimited) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

// VerifyEmail 校验验证码并完成绑定：成功后更新该 wallet 的 User.Email。
// R1 防越权：归属 wallet 唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body。
// 验证码错误/过期/锁死 → 400「验证码错误或已过期」。
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authWallet := c.GetString(middleware.CtxWallet)
	if err := h.emailService.ConfirmBind(authWallet, req.Email, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "邮箱绑定成功"})
}

type RegisterRequest struct {
	Wallet  string `json:"wallet" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	TokenID uint   `json:"token_id"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.Register(req.Wallet, req.Email, req.TokenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *Handler) GetUser(c *gin.Context) {
	wallet := c.Param("wallet")
	user, err := h.userService.GetUser(wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetOperators(c *gin.Context) {
	operators, err := h.operatorService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, operators)
}

type ActivateServiceRequest struct {
	Wallet        string `json:"wallet" binding:"required"`
	OperatorID    uint   `json:"operator_id" binding:"required"`
	VirtualNumber string `json:"virtual_number" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

func (h *Handler) ActivateService(c *gin.Context) {
	var req ActivateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// R1 防越权：操作主体唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body.wallet
	// （攻击者用自签名过闸、body 填受害者地址操作他人服务的横向越权面在此收口）。
	authWallet := c.GetString(middleware.CtxWallet)
	err := h.userServiceService.Activate(authWallet, req.OperatorID, req.VirtualNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service activated successfully"})
}

func (h *Handler) DeactivateService(c *gin.Context) {
	// R1 防越权：停用主体唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body
	// （否则攻击者可即时停用他人服务）。body 不再读 wallet。
	authWallet := c.GetString(middleware.CtxWallet)
	err := h.userServiceService.Deactivate(authWallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deactivated successfully"})
}

type DepositRequest struct {
	Wallet string `json:"wallet" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

func (h *Handler) Deposit(c *gin.Context) {
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// R1 防越权：充值意向归属唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body.wallet。
	authWallet := c.GetString(middleware.CtxWallet)
	err := h.depositService.Deposit(authWallet, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

func (h *Handler) GetDeposit(c *gin.Context) {
	wallet := c.Param("wallet")
	amount, err := h.depositService.GetDepositAmount(wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

// Withdraw 接收提现 pending 意向（B3 降级，design §4.3/§6.6）。
// ⚠️ 不凭前端 tx_hash 写记账：提现终态唯一由 event_sync 监听 DepositWithdrawn 回填（等 K 块确认）。
// 这里只落 Status=pending 的意向（不计入余额）。鉴权由 WalletAuth 中间件保证（钱包签名 ecrecover）。
func (h *Handler) Withdraw(c *gin.Context) {
	var req struct {
		TxHash string `json:"tx_hash"`
	}
	_ = c.ShouldBindJSON(&req) // body 仅含可选 tx_hash；wallet 不再从 body 取
	// R1 防越权：提现意向归属唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body.wallet。
	authWallet := c.GetString(middleware.CtxWallet)
	if err := h.depositService.RecordPendingWithdraw(authWallet, req.TxHash); err != nil {
		log.Printf("Withdraw intent failed wallet=%s tx=%s err=%v", authWallet, req.TxHash, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Withdraw intent recorded (pending) wallet=%s tx=%s", authWallet, req.TxHash)
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal intent recorded (pending; confirmed by on-chain event)"})
}

func (h *Handler) GetDepositHistory(c *gin.Context) {
	wallet := c.Param("wallet")
	records, err := h.depositService.GetHistory(wallet)
	if err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *Handler) GetBills(c *gin.Context) {
	wallet := c.Param("wallet")
	bills, err := h.billingService.GetBills(wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bills)
}

// PayBill 接收用户「我已发起支付」意向（B2 降级，design §4.3/§6.6）。
// ⚠️ 只写 Bill.PayIntentTxHash，绝不置 IsPaid——IsPaid 终态唯一由 event_sync 监听 BillPaid 回填。
// 前端不能据本端点 200 就显示「已付」，应轮询 GET /api/bills/:wallet 看 is_paid。
// 鉴权由 WalletAuth 中间件保证（钱包签名 ecrecover 绑 wallet，防冒充他人提交意向）。
func (h *Handler) PayBill(c *gin.Context) {
	var req struct {
		BillID uint   `json:"bill_id" binding:"required"`
		TxHash string `json:"tx_hash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// R1 防越权：意向归属唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body.wallet。
	// RecordPayIntent→SetPayIntent 以 authWallet 解析的 userID 限定，确保 bill 属于已鉴权钱包
	// （攻击者替他人账单写意向被 user_id 不匹配拒绝）。
	authWallet := c.GetString(middleware.CtxWallet)
	if err := h.billingService.RecordPayIntent(authWallet, req.BillID, req.TxHash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment intent recorded (pending; confirmed by on-chain BillPaid event)"})
}

// SimClaimRequest 「流量卡销毁兑换 SIM」claim 意向（新玩法）。
// tokenIds 为前端勾选销毁的流量卡 NFT（每张=1 天），days = len(tokenIds)。
// txHash 为链上 redeemForSim 交易哈希（意向追溯；SIM confirmed 终态唯一由 event_sync 监听 SimRedeemed 回填）。
type SimClaimRequest struct {
	Destination string   `json:"destination" binding:"required"`
	Recipient   string   `json:"recipient" binding:"required"`
	AddressLine string   `json:"addressLine" binding:"required"`
	TokenIDs    []uint64 `json:"tokenIds" binding:"required,min=1"`
	TxHash      string   `json:"txHash"`
}

// ClaimSim 接收 SIM 兑换 pending 意向（WalletAuth 保护，action "sim/claim"）。
// ⚠️ 只落 Status=pending 意向——SIM confirmed 终态唯一由 event_sync 监听链上 SimRedeemed 等 K 块回填。
// R1 防越权：归属 wallet 唯一取自 WalletAuth 已验证钱包（CtxWallet），不信任 body。
func (h *Handler) ClaimSim(c *gin.Context) {
	var req SimClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authWallet := c.GetString(middleware.CtxWallet)
	sim, err := h.simService.Claim(authWallet, req.Destination, req.Recipient, req.AddressLine, uint(len(req.TokenIDs)), req.TxHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "SIM redemption intent recorded (pending; confirmed by on-chain SimRedeemed event)",
		"sim":     sim,
	})
}

// GetSims 公开读：返回该用户的 SIM 列表（camelCase JSON）。
func (h *Handler) GetSims(c *gin.Context) {
	wallet := c.Param("wallet")
	sims, err := h.simService.GetByWallet(wallet)
	if err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, sims)
}

type GenerateVirtualNumberRequest struct {
	CountryCode string `json:"country_code" binding:"required"`
}

func (h *Handler) GenerateVirtualNumber(c *gin.Context) {
	var req GenerateVirtualNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	number, password, err := h.virtualGen.Generate(req.CountryCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"virtual_number": number,
		"password":       password,
	})
}

func (h *Handler) GetUsage(c *gin.Context) {
	wallet := c.Param("wallet")
	dataUsage, callUsage, signature, err := h.usageService.QueryUsage(wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data_usage": dataUsage,
		"call_usage": callUsage,
		"signature":  signature,
	})
}

// TriggerMonthlyBill 平台触发月度结算（design §6.5/§6.6，AdminAuth 保护）。
// 两段：① 始终落 DB Bill（usageService.TriggerMonthlyBill，账单记录）；
// ② 若结算编排器已装配（owner key/链就绪），调 oracleV2.SettleMonthlyOnChain(ctx, month)
//
//	组批上链并返回 SettlementSummary（分批摘要：confirmed/failed/skipped/blocked/pending_review + txHashes）。
//
// 未装配编排器（无 owner key/未配链）→ 仅返回 DB bill count + settlement=null（降级，不报错）。
// month 取请求体 month（"YYYY-MM"），缺省用当前月。
func (h *Handler) TriggerMonthlyBill(c *gin.Context) {
	var req struct {
		Month string `json:"month"`
	}
	_ = c.ShouldBindJSON(&req) // month 可选，绑定失败（空 body）不阻断
	month := req.Month
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	count, err := h.usageService.TriggerMonthlyBill()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summary, serr := h.oracleV2.SettleMonthlyOnChain(c.Request.Context(), month)
	if serr != nil {
		// 编排器未装配（无 owner key/未配链）属预期降级：返回 DB count + 降级说明，不算错误。
		log.Printf("INFO: 月度上链结算降级（仅落 DB bill）：%v", serr)
		c.JSON(http.StatusOK, gin.H{
			"message":    "Monthly bills created (on-chain settlement skipped: not configured)",
			"count":      count,
			"month":      month,
			"settlement": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Monthly bills created and settled on-chain",
		"count":      count,
		"month":      month,
		"settlement": summary,
	})
}

func (h *Handler) GetCountryList(c *gin.Context) {
	list := h.virtualGen.GetCountryList()
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetUserService(c *gin.Context) {
	wallet := c.Param("wallet")
	service, err := h.userServiceService.GetUserService(wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active service"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// SubmitUsageRequest 用量上报（B4 范围校验，design §6.6/§7.0②）：
// data_usage(MB)/call_usage(min) gin binding `max=` 上界与 PricingService 的 MaxDataMB/MaxCallMin 同源
// （services.MaxDataMB=1_000_000、services.MaxCallMin=100_000），超界 → 400，防天文账单在入口拦截。
// CallUsage 不加 required（合法 0 通话量不应被拒），用 gte=0 仅约束非负。
type SubmitUsageRequest struct {
	Wallet     string `json:"wallet" binding:"required"`
	OperatorID uint   `json:"operator_id" binding:"required"`
	DataUsage  uint64 `json:"data_usage" binding:"max=1000000"`
	CallUsage  uint64 `json:"call_usage" binding:"max=100000"`
	Signature  string `json:"signature" binding:"required"`
}

func (h *Handler) SubmitUsage(c *gin.Context) {
	var req SubmitUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.oracleService.SubmitUsage(req.Wallet, req.OperatorID, req.DataUsage, req.CallUsage, req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usage submitted"})
}

func (h *Handler) SendNotification(c *gin.Context) {
	var req struct {
		Wallet string `json:"wallet" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bills, err := h.billingService.GetUnpaidBills(req.Wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.notificationService.SendBillEmail(req.Wallet, bills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}
