package services

import (
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

func (s *BillingService) MarkAsPaid(billID uint, txHash string) error {
	return s.billRepo.MarkAsPaid(billID, txHash)
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
}

func NewOracleService(userServiceRepo *repository.UserServiceRepository, billRepo *repository.BillRepository) *OracleService {
	return &OracleService{
		userServiceRepo: userServiceRepo,
		billRepo:        billRepo,
	}
}

func (s *OracleService) SubmitUsage(userWallet string, operatorID uint, dataUsage, callUsage uint64, signature string) error {
	return nil
}

type NotificationService struct {
	userRepo *repository.UserRepository
}

func NewNotificationService(userRepo *repository.UserRepository) *NotificationService {
	return &NotificationService{userRepo: userRepo}
}

func (s *NotificationService) SendBillEmail(wallet string, bills []models.Bill) error {
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

func (s *DepositService) Deposit(wallet, amount string) error {
	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return err
	}
	deposit := &models.Deposit{
		UserID: user.ID,
		Amount: amount,
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
	return s.userServiceRepo.FindByUserID(user.ID)
}

var ErrServiceAlreadyActive = &ServiceError{"Service already active"}
