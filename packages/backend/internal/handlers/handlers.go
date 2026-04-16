package handlers

import (
	"log"
	"net/http"

	"linkworld-backend/internal/services"

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
	}
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

	err := h.userServiceService.Activate(req.Wallet, req.OperatorID, req.VirtualNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service activated successfully"})
}

func (h *Handler) DeactivateService(c *gin.Context) {
	var req struct {
		Wallet string `json:"wallet" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userServiceService.Deactivate(req.Wallet)
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

	err := h.depositService.Deposit(req.Wallet, req.Amount)
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

func (h *Handler) Withdraw(c *gin.Context) {
	var req struct {
		Wallet string `json:"wallet" binding:"required"`
		TxHash string `json:"tx_hash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Withdrawal is completed on-chain; backend only acknowledges the event
	log.Printf("Withdraw acknowledged for wallet=%s tx=%s", req.Wallet, req.TxHash)
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal acknowledged"})
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

func (h *Handler) PayBill(c *gin.Context) {
	var req struct {
		Wallet string `json:"wallet" binding:"required"`
		BillID uint   `json:"bill_id" binding:"required"`
		TxHash string `json:"tx_hash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.billingService.MarkAsPaid(req.BillID, req.TxHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bill paid successfully"})
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

func (h *Handler) TriggerMonthlyBill(c *gin.Context) {
	count, err := h.usageService.TriggerMonthlyBill()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Monthly bills created", "count": count})
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
