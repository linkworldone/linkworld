package handlers

import (
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
}

func NewHandler(
	userService *services.UserService,
	operatorService *services.OperatorService,
	billingService *services.BillingService,
	oracleService *services.OracleService,
	notificationService *services.NotificationService,
	depositService *services.DepositService,
	userServiceService *services.UserServiceService,
) *Handler {
	return &Handler{
		userService:         userService,
		operatorService:     operatorService,
		billingService:      billingService,
		oracleService:       oracleService,
		notificationService: notificationService,
		depositService:      depositService,
		userServiceService:  userServiceService,
	}
}

type RegisterRequest struct {
	Wallet string `json:"wallet" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.Register(req.Wallet, req.Email, 0)
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
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Withdraw request submitted, please wait for blockchain confirmation"})
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
	c.JSON(http.StatusOK, gin.H{"message": "Bill paid"})
}
