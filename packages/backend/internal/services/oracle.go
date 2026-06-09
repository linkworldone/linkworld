package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
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
		"US": "+1",
		"GB": "+44",
		"FR": "+33",
		"RU": "+7",
		"JP": "+81",
		"VN": "+84",
		"LA": "+856",
		"KH": "+855",
		"TH": "+66",
		"MY": "+60",
		"PH": "+63",
	}

	CountryLengths = map[string]int{
		"US": 10, "GB": 10, "FR": 9, "RU": 10, "JP": 10,
		"VN": 9, "LA": 8, "KH": 8, "TH": 9, "MY": 9, "PH": 10,
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
		"US": "United States", "GB": "United Kingdom", "FR": "France", "RU": "Russia",
		"JP": "Japan", "VN": "Vietnam", "LA": "Laos", "KH": "Cambodia",
		"TH": "Thailand", "MY": "Malaysia", "PH": "Philippines",
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

// OperatorAPI 提供运营商用量（usage）。
//
// 注（design §3.2 / arch-review CEO 验收语言）：本轮 usage 仍为模拟（真实运营商 API 留后续）；
// 计价金额已改为确定可复算（PricingService），但「计价引擎正确 ≠ 计费正确」。
// 原 GetBill（返回 rand.Intn 随机金额 + 拍脑袋 2.5% fee）已废弃删除——金额唯一由 PricingService 计算（design §4.2）。
type OperatorAPI interface {
	GetUsage(userID, operatorID uint) (uint64, uint64, error)
}

type OperatorAPISimulator struct{}

func NewOperatorAPISimulator() *OperatorAPISimulator {
	return &OperatorAPISimulator{}
}

// GetUsage 模拟运营商用量（单位锁：dataUsage=MB、callUsage=minute，design §4.4）。
// 本轮模拟值刻意限制在 usage 上界（MaxDataMB/MaxCallMin）内，避免触发 L1 上界拒绝。
func (s *OperatorAPISimulator) GetUsage(userID, operatorID uint) (uint64, uint64, error) {
	dataUsage := uint64(rand.Intn(10000) + 100) // MB，远低于 MaxDataMB
	callUsage := uint64(rand.Intn(500))         // minute，远低于 MaxCallMin
	return dataUsage, callUsage, nil
}

type OracleServiceV2 struct {
	operatorAPI     OperatorAPI
	pricing         *PricingService
	userRepo        *repository.UserRepository
	userServiceRepo *repository.UserServiceRepository
	billRepo        *repository.BillRepository
	usageRepo       *repository.UsageDataRepository
}

func NewOracleServiceV2(
	operatorAPI OperatorAPI,
	pricing *PricingService,
	userRepo *repository.UserRepository,
	userServiceRepo *repository.UserServiceRepository,
	billRepo *repository.BillRepository,
	usageRepo *repository.UsageDataRepository,
) *OracleServiceV2 {
	return &OracleServiceV2{
		operatorAPI:     operatorAPI,
		pricing:         pricing,
		userRepo:        userRepo,
		userServiceRepo: userServiceRepo,
		billRepo:        billRepo,
		usageRepo:       usageRepo,
	}
}

func (s *OracleServiceV2) FetchAndCreateBills() (int, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return 0, err
	}

	now := time.Now()

	billCount := 0
	for _, user := range users {
		userService, err := s.userServiceRepo.GetActiveByUserID(user.ID)
		if err != nil || userService == nil || !userService.IsActive {
			continue
		}

		// 用量（模拟）→ 计价（确定可复算，替换原 rand.Intn 随机金额，design §4.2）。
		// 单位锁：dataMB=MB、callMin=minute。
		dataMB, callMin, err := s.operatorAPI.GetUsage(user.ID, userService.OperatorID)
		if err != nil {
			continue
		}
		// L1 计价 + usage 上界 + 单 bill 硬上限（design §7.0①）。超界/超额：拒绝该条 + 告警，不出账。
		amount6, err := s.pricing.Price(dataMB, callMin, userService.OperatorID)
		if err != nil {
			log.Printf("WARN: 计价拒绝(L1)：user=%d op=%d dataMB=%d callMin=%d err=%v（跳过该 bill）", user.ID, userService.OperatorID, dataMB, callMin, err)
			continue
		}
		// amount6 == 0（低于最小出账额）→ 跳过（合约侧 amount>0 才 createBill，design §4.1）。
		if amount6.Sign() == 0 {
			continue
		}

		bill := &models.Bill{
			UserID:               user.ID,
			OperatorID:           userService.OperatorID,
			Amount:               amount6.String(), // 6 位最小单位字符串（design §4.4）
			PlatformFee:          "0",              // 平台费由合约 FeeManager 计算并 emit，后端不再拍脑袋 2.5%
			TrafficCardDeduction: "0",              // 本轮恒 "0"（design §4.3，applyTrafficCardToBill 桩）
			IsPaid:               false,
			CreatedAt:            now,
		}
		s.billRepo.Create(bill)
		billCount++
	}

	return billCount, nil
}

func (s *OracleServiceV2) FetchUsage(userWallet string) (uint64, uint64, uint, error) {
	user, err := s.userRepo.FindByWallet(userWallet)
	if err != nil {
		return 0, 0, 0, err
	}

	userService, err := s.userServiceRepo.GetActiveByUserID(user.ID)
	if err != nil || userService == nil || !userService.IsActive {
		return 0, 0, 0, nil
	}

	dataUsage, callUsage, err := s.operatorAPI.GetUsage(user.ID, userService.OperatorID)
	if err != nil {
		return 0, 0, 0, err
	}

	return dataUsage, callUsage, userService.OperatorID, nil
}

func (s *OracleServiceV2) SignData(userID uint, operatorID uint, dataUsage, callUsage uint64) string {
	data := fmt.Sprintf("%d-%d-%d-%d", userID, operatorID, dataUsage, callUsage)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

type UsageService struct {
	oracleService   *OracleServiceV2
	usageRepo       *repository.UsageDataRepository
	userRepo        *repository.UserRepository
	userServiceRepo *repository.UserServiceRepository
}

func NewUsageService(oracleService *OracleServiceV2, usageRepo *repository.UsageDataRepository, userRepo *repository.UserRepository, userServiceRepo *repository.UserServiceRepository) *UsageService {
	return &UsageService{
		oracleService:   oracleService,
		usageRepo:       usageRepo,
		userRepo:        userRepo,
		userServiceRepo: userServiceRepo,
	}
}

func (s *UsageService) QueryUsage(wallet string) (uint64, uint64, string, error) {
	dataUsage, callUsage, operatorID, err := s.oracleService.FetchUsage(wallet)
	if err != nil {
		return 0, 0, "", err
	}

	user, err := s.userRepo.FindByWallet(wallet)
	if err != nil {
		return 0, 0, "", err
	}

	signature := s.oracleService.SignData(user.ID, uint(operatorID), dataUsage, callUsage)

	usageData := &models.UsageData{
		UserID:     user.ID,
		OperatorID: uint(operatorID),
		DataUsage:  dataUsage,
		CallUsage:  callUsage,
		Timestamp:  time.Now(),
		Signature:  signature,
	}
	s.usageRepo.Create(usageData)

	return dataUsage, callUsage, signature, nil
}

func (s *UsageService) TriggerMonthlyBill() (int, error) {
	count, err := s.oracleService.FetchAndCreateBills()
	if err != nil {
		return 0, err
	}
	return count, nil
}