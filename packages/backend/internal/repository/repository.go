package repository

import (
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

func (r *BillRepository) MarkAsPaid(id uint) error {
	return r.db.Model(&models.Bill{}).Where("id = ?", id).Update("is_paid", true).Error
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
	err := r.db.Where("user_id = ?", userID).Find(&deposits).Error
	return deposits, err
}

func (r *DepositRepository) GetTotalByUserID(userID uint) (string, error) {
	var result struct {
		Total string
	}
	err := r.db.Model(&models.Deposit{}).Select("COALESCE(SUM(amount), '0') as total").Where("user_id = ?", userID).Scan(&result).Error
	return result.Total, err
}
