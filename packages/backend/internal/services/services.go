package services

import (
	"fmt"
	"log"
	"time"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(wallet, email string, tokenID uint) error {
	if s.repo.Exists(wallet) {
		return ErrAlreadyRegistered
	}
	user := &models.User{
		WalletAddr: wallet,
		Email:      email,
		TokenID:    tokenID,
		IsActive:   true,
	}
	return s.repo.Create(user)
}

func (s *UserService) GetUser(wallet string) (*models.User, error) {
	return s.repo.FindByWallet(wallet)
}

var ErrAlreadyRegistered = &ServiceError{"User already registered"}

type ServiceError struct {
	msg string
}

func (e *ServiceError) Error() string {
	return e.msg
}

type OperatorService struct {
	repo *repository.OperatorRepository
}

func NewOperatorService(repo *repository.OperatorRepository) *OperatorService {
	return &OperatorService{repo: repo}
}

func (s *OperatorService) GetAll() ([]models.Operator, error) {
	return s.repo.FindAll()
}

func (s *OperatorService) GetByID(id uint) (*models.Operator, error) {
	return s.repo.FindByID(id)
}

type BillingService struct {
	billRepo *repository.BillRepository
	userRepo *repository.UserRepository
}

func NewBillingService(billRepo *repository.BillRepository, userRepo *repository.UserRepository) *BillingService {
	return &BillingService{
		billRepo: billRepo,
		userRepo: userRepo,
	}
}

func (s *BillingService) GetBills(wallet string) ([]models.Bill, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return s.billRepo.FindByUserID(user.ID)
}

// MarkAsPaid 唯一由 event_sync(T4) 内部调用（监听 BillPaid 确认后置 IsPaid）。
// ⚠️ 绝不暴露给 HTTP handler（B2，design §4.3/§6.6）：HTTP 侧只走 RecordPayIntent 写 pending 意向。
func (s *BillingService) MarkAsPaid(billID uint, txHash string) error {
	return s.billRepo.MarkAsPaid(billID, txHash)
}

// RecordPayIntent 记录用户「我已发起支付」意向（B2 降级，design §4.3/§6.6）：
// 只写 Bill.PayIntentTxHash，绝不置 IsPaid——IsPaid 终态唯一由 event_sync 监听 BillPaid 回填。
// 校验 bill 属于该 wallet（防越权改他人账单意向）。
func (s *BillingService) RecordPayIntent(wallet string, billID uint, txHash string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	return s.billRepo.SetPayIntent(user.ID, billID, txHash)
}

func (s *BillingService) GetUnpaidBills(wallet string) ([]models.Bill, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return s.billRepo.FindUnpaidByUserID(user.ID)
}

type OracleService struct {
	userServiceRepo *repository.UserServiceRepository
	billRepo        *repository.BillRepository
	usageRepo       *repository.UsageDataRepository
	userRepo        *repository.UserRepository
}

func NewOracleService(userServiceRepo *repository.UserServiceRepository, billRepo *repository.BillRepository, usageRepo *repository.UsageDataRepository, userRepo *repository.UserRepository) *OracleService {
	return &OracleService{
		userServiceRepo: userServiceRepo,
		billRepo:        billRepo,
		usageRepo:       usageRepo,
		userRepo:        userRepo,
	}
}

func (s *OracleService) SubmitUsage(userWallet string, operatorID uint, dataUsage, callUsage uint64, signature string) error {
	user, err := s.userRepo.FindByWallet(userWallet)
	if err != nil {
		return err
	}

	usageData := &models.UsageData{
		UserID:     user.ID,
		OperatorID: operatorID,
		DataUsage:  dataUsage,
		CallUsage:  callUsage,
		Signature:  signature,
		Timestamp:  time.Now(),
	}
	return s.usageRepo.Create(usageData)
}

type NotificationService struct {
	userRepo *repository.UserRepository
}

func NewNotificationService(userRepo *repository.UserRepository) *NotificationService {
	return &NotificationService{userRepo: userRepo}
}

func (s *NotificationService) SendBillEmail(wallet string, bills []models.Bill) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	if user.Email == "" {
		return fmt.Errorf("no email registered for user")
	}
	// TODO: Integrate with email service (SMTP, SendGrid, etc.)
	// For now, log the intent
	log.Printf("Sending %d bills to email %s", len(bills), user.Email)
	return nil
}

type DepositService struct {
	depositRepo *repository.DepositRepository
	userRepo    *repository.UserRepository
}

func NewDepositService(depositRepo *repository.DepositRepository, userRepo *repository.UserRepository) *DepositService {
	return &DepositService{
		depositRepo: depositRepo,
		userRepo:    userRepo,
	}
}

// Deposit 记录用户充值 pending 意向（design §4.3/§4.4 状态机，B3/B5）。
// ⚠️ 写 Status=pending（不计入余额）——余额唯一由 event_sync 监听 DepositMade 等 K 块置 confirmed。
// 旧实现写空 status（视为历史 confirmed）会把未上链意向直接计入余额，已收口为 pending。
func (s *DepositService) Deposit(wallet, amount string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	deposit := &models.Deposit{
		UserID: user.ID,
		Amount: amount,
		Type:   "deposit",
		Status: models.DepositStatusPending,
	}
	return s.depositRepo.Create(deposit)
}

func (s *DepositService) GetDepositAmount(wallet string) (string, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return "", err
	}
	return s.depositRepo.GetTotalByUserID(user.ID)
}

// RecordPendingWithdraw 记录用户提现 pending 意向（B3 降级，design §4.3/§6.6）。
// ⚠️ 不凭前端 txHash 写记账（旧 RecordWithdraw 凭 txHash 直接写 withdraw 记录已废弃，伪造记账面）：
// 提现终态（实际扣减本金/利息）唯一由 event_sync 监听 DepositWithdrawn(user,principal,interest) 回填。
// 这里只落一条 Status=pending 的意向记录（不计入余额，GetTotalByUserID 只认 confirmed）。
// txHash 仅作意向追溯，不作记账依据；金额留空（真实金额以链上事件为准）。
func (s *DepositService) RecordPendingWithdraw(wallet, txHash string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	deposit := &models.Deposit{
		UserID: user.ID,
		Amount: "0", // 意向占位；真实提现金额唯一由 DepositWithdrawn 事件回填
		Type:   "withdraw",
		TxHash: txHash,
		Status: models.DepositStatusPending,
	}
	return s.depositRepo.Create(deposit)
}

func (s *DepositService) GetHistory(wallet string) ([]models.Deposit, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return s.depositRepo.FindByUserID(user.ID)
}

type UserServiceService struct {
	userServiceRepo *repository.UserServiceRepository
	operatorRepo    *repository.OperatorRepository
	userRepo        *repository.UserRepository
}

func NewUserServiceService(userServiceRepo *repository.UserServiceRepository, operatorRepo *repository.OperatorRepository, userRepo *repository.UserRepository) *UserServiceService {
	return &UserServiceService{
		userServiceRepo: userServiceRepo,
		operatorRepo:    operatorRepo,
		userRepo:        userRepo,
	}
}

func (s *UserServiceService) Activate(wallet string, operatorID uint, virtualNumber, password string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	_, err = s.operatorRepo.FindByID(operatorID)
	if err != nil {
		return err
	}
	existing, _ := s.userServiceRepo.FindByUserID(user.ID)
	if existing != nil && existing.IsActive {
		return ErrServiceAlreadyActive
	}
	userService := &models.UserService{
		UserID:        user.ID,
		OperatorID:    operatorID,
		VirtualNumber: virtualNumber,
		Password:      password,
		IsActive:      true,
	}
	if existing != nil {
		existing.OperatorID = operatorID
		existing.VirtualNumber = virtualNumber
		existing.Password = password
		existing.IsActive = true
		return s.userServiceRepo.Update(existing)
	}
	return s.userServiceRepo.Create(userService)
}

func (s *UserServiceService) Deactivate(wallet string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	return s.userServiceRepo.Deactivate(user.ID)
}

func (s *UserServiceService) GetUserService(wallet string) (*models.UserService, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return s.userServiceRepo.GetActiveByUserID(user.ID)
}

var ErrServiceAlreadyActive = &ServiceError{"Service already active"}

// SimService 「流量卡销毁兑换 SIM」业务（新玩法，design 新增）。
// 用户在前端勾选多张流量卡 NFT 一次性 redeemForSim 销毁，兑换到一张无限流量 SIM。
// 链上只销毁卡 + emit SimRedeemed；SIM 身份/天数/收件信息全记后端 DB。
type SimService struct {
	simRepo  *repository.SimRepository
	userRepo *repository.UserRepository
}

func NewSimService(simRepo *repository.SimRepository, userRepo *repository.UserRepository) *SimService {
	return &SimService{simRepo: simRepo, userRepo: userRepo}
}

// Claim 记录用户 SIM 兑换 pending 意向（同押金两阶段状态机）。
// ⚠️ 写 Status=pending——SIM 终态(confirmed)唯一由 event_sync 监听 SimRedeemed 等 K 块回填。
// days = 销毁卡数（len(tokenIds)），归属 wallet 由调用方从 WalletAuth CtxWallet 取（不信 body）。
func (s *SimService) Claim(wallet, destination, recipient, addressLine string, days uint, txHash string) (*models.Sim, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	sim := &models.Sim{
		UserID:      user.ID,
		Days:        days,
		Destination: destination,
		Recipient:   recipient,
		AddressLine: addressLine,
		TxHash:      txHash,
		Status:      models.SimStatusPending,
	}
	if err := s.simRepo.Create(sim); err != nil {
		return nil, err
	}
	return sim, nil
}

// GetByWallet 返回该用户的 SIM 列表。
func (s *SimService) GetByWallet(wallet string) ([]models.Sim, error) {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return s.simRepo.FindByUserID(user.ID)
}
