package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
)

type VirtualNumberGenerator struct{}

func NewVirtualNumberGenerator() *VirtualNumberGenerator {
	return &VirtualNumberGenerator{}
}

var (
	CountryCodes = map[string]string{
		"US": "+1",   // 美国
		"GB": "+44",  // 英国
		"FR": "+33",  // 法国
		"RU": "+7",   // 俄罗斯
		"JP": "+81",  // 日本
		"VN": "+84",  // 越南
		"LA": "+856", // 老挝
		"KH": "+855", // 柬埔寨
		"TH": "+66",  // 泰国
		"MY": "+60",  // 马来西亚
		"PH": "+63",  // 菲律宾
	}

	CountryLengths = map[string]int{
		"US": 10,
		"GB": 10,
		"FR": 9,
		"RU": 10,
		"JP": 10,
		"VN": 9,
		"LA": 8,
		"KH": 8,
		"TH": 9,
		"MY": 9,
		"PH": 10,
	}

	CountryPrefixes = map[string][]string{
		"US": {"201", "305", "312", "415", "502", "510", "602", "617", "702", "714", "718", "805", "818", "860", "901", "908"},
		"GB": {"07", "020", "021", "022", "023", "024", "028", "029"},
		"FR": {"06", "07"},
		"RU": {"9", "8"},
		"JP": {"70", "80", "90"},
		"VN": {"84", "91", "94", "98"},
		"LA": {"20", "21"},
		"KH": {"10", "11", "12", "60", "66", "71", "76", "81", "86", "92", "96", "99"},
		"TH": {"06", "08", "09"},
		"MY": {"10", "11", "12"},
		"PH": {"09"},
	}

	CountryNames = map[string]string{
		"US": "United States",
		"GB": "United Kingdom",
		"FR": "France",
		"RU": "Russia",
		"JP": "Japan",
		"VN": "Vietnam",
		"LA": "Laos",
		"KH": "Cambodia",
		"TH": "Thailand",
		"MY": "Malaysia",
		"PH": "Philippines",
	}
)

func (g *VirtualNumberGenerator) Generate(countryCode string) (string, string, error) {
	code, ok := CountryCodes[countryCode]
	if !ok {
		return "", "", fmt.Errorf("unsupported country: %s", countryCode)
	}

	length := CountryLengths[countryCode]
	prefixes := CountryPrefixes[countryCode]

	prefix := prefixes[rand.Intn(len(prefixes))]
	if code == "+1" || code == "+44" || code == "+33" || code == "+7" {
		prefix = prefixes[rand.Intn(len(prefixes))]
	}

	number := g.generateDigits(length - len(prefix))
	fullNumber := code + prefix + number

	password := g.generatePassword(8)

	return fullNumber, password, nil
}

func (g *VirtualNumberGenerator) GetCountryList() []map[string]string {
	list := make([]map[string]string, 0, len(CountryCodes))
	for code, name := range CountryNames {
		list = append(list, map[string]string{
			"code":   code,
			"name":   name,
			"prefix": CountryCodes[code],
		})
	}
	return list
}

func (g *VirtualNumberGenerator) generateDigits(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

func (g *VirtualNumberGenerator) generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

type OperatorAPI interface {
	GetUsage(userID, operatorID uint) (uint64, uint64, error)
	GetBill(userID, operatorID uint, month time.Time) (string, string, error)
}

type OperatorAPISimulator struct{}

func NewOperatorAPISimulator() *OperatorAPISimulator {
	return &OperatorAPISimulator{}
}

func (s *OperatorAPISimulator) GetUsage(userID, operatorID uint) (uint64, uint64, error) {
	dataUsage := uint64(rand.Intn(10000) + 100)
	callUsage := uint64(rand.Intn(500))
	return dataUsage, callUsage, nil
}

func (s *OperatorAPISimulator) GetBill(userID, operatorID uint, month time.Time) (string, string, error) {
	baseAmount := uint64(rand.Intn(5000) + 500)
	fee := baseAmount * 25 / 1000
	return fmt.Sprintf("%d", baseAmount), fmt.Sprintf("%d", fee), nil
}

type OracleServiceV2 struct {
	operatorAPI OperatorAPI
	userRepo    *repository.UserRepository
	billRepo    *repository.BillRepository
	usageRepo   *repository.UsageDataRepository
}

func NewOracleServiceV2(
	operatorAPI OperatorAPI,
	userRepo *repository.UserRepository,
	billRepo *repository.BillRepository,
	usageRepo *repository.UsageDataRepository,
) *OracleServiceV2 {
	return &OracleServiceV2{
		operatorAPI: operatorAPI,
		userRepo:    userRepo,
		billRepo:    billRepo,
		usageRepo:   usageRepo,
	}
}

func (s *OracleServiceV2) FetchAndCreateBills() error {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return err
	}

	now := time.Now()
	month := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())

	for _, user := range users {
		services, err := s.usageRepo.FindByUserID(user.ID)
		if err != nil || len(services) == 0 {
			continue
		}

		for _, service := range services {
			amount, fee, err := s.operatorAPI.GetBill(user.ID, service.OperatorID, month)
			if err != nil {
				continue
			}

			bill := &models.Bill{
				UserID:      user.ID,
				OperatorID:  service.OperatorID,
				Amount:      amount,
				PlatformFee: fee,
				IsPaid:      false,
				CreatedAt:   now,
			}
			s.billRepo.Create(bill)
		}
	}

	return nil
}

func (s *OracleServiceV2) FetchUsage(userWallet string) (uint64, uint64, error) {
	user, err := s.userRepo.FindByWallet(userWallet)
	if err != nil {
		return 0, 0, err
	}

	services, err := s.usageRepo.FindByUserID(user.ID)
	if err != nil || len(services) == 0 {
		return 0, 0, nil
	}

	service := services[0]
	return s.operatorAPI.GetUsage(user.ID, service.OperatorID)
}

func (s *OracleServiceV2) SignData(userID uint, operatorID uint, dataUsage, callUsage uint64) string {
	data := fmt.Sprintf("%d-%d-%d-%d", userID, operatorID, dataUsage, callUsage)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

type UsageService struct {
	oracleService *OracleServiceV2
	usageRepo     *repository.UsageDataRepository
	userRepo      *repository.UserRepository
}

func NewUsageService(oracleService *OracleServiceV2, usageRepo *repository.UsageDataRepository, userRepo *repository.UserRepository) *UsageService {
	return &UsageService{
		oracleService: oracleService,
		usageRepo:     usageRepo,
		userRepo:      userRepo,
	}
}

func (s *UsageService) QueryUsage(wallet string) (uint64, uint64, string, error) {
	dataUsage, callUsage, err := s.oracleService.FetchUsage(wallet)
	if err != nil {
		return 0, 0, "", err
	}

	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return 0, 0, "", err
	}

	signature := s.oracleService.SignData(user.ID, 1, dataUsage, callUsage)

	usageData := &models.UsageData{
		UserID:     user.ID,
		OperatorID: 1,
		DataUsage:  dataUsage,
		CallUsage:  callUsage,
		Timestamp:  time.Now(),
		Signature:  signature,
	}
	s.usageRepo.Create(usageData)

	return dataUsage, callUsage, signature, nil
}

func (s *UsageService) TriggerMonthlyBill() (int, error) {
	err := s.oracleService.FetchAndCreateBills()
	if err != nil {
		return 0, err
	}
	return 1, nil
}
