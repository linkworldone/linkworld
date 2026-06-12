package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
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
	// settlement 是 T6 组批结算编排器（可选，nil 时仅走 DB Bill 落库不上链——无 owner key/未配链时）。
	settlement *SettlementOrchestrator
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

// SetSettlementOrchestrator 注入 T6 组批结算编排器（main.go 在 owner key/链配置就绪时装配）。
func (s *OracleServiceV2) SetSettlementOrchestrator(o *SettlementOrchestrator) {
	s.settlement = o
}

func (s *OracleServiceV2) FetchAndCreateBills() (int, error) {
	items, _, err := s.collectSettlementItems()
	if err != nil {
		return 0, err
	}

	now := time.Now()
	billCount := 0
	for _, it := range items {
		bill := &models.Bill{
			UserID:               it.userID,
			OperatorID:           it.OperatorID,
			Amount:               it.Amount6.String(), // 6 位最小单位字符串（design §4.4）
			PlatformFee:          "0",                 // 平台费由合约 FeeManager 计算并 emit，后端不再拍脑袋 2.5%
			TrafficCardDeduction: "0",                 // 本轮恒 "0"（design §4.3，applyTrafficCardToBill 桩）
			IsPaid:               false,
			CreatedAt:            now,
		}
		s.billRepo.Create(bill)
		billCount++
	}

	return billCount, nil
}

// settlementItem 是 collectSettlementItems 内部产物：携带 DB userID（落 Bill 用）+ 链上结算项。
type settlementItem struct {
	userID uint
	SettlementItem
}

// collectSettlementItems 遍历激活服务用户，逐 user 走 L1 计价（PricingService + usage 上界 + 单 bill 硬上限），
// 产出待结算项（含链地址 + ChainOperatorID 映射 + amount6）。L1 拒绝/零额项跳过（不入数组）。
// 返回 (items, skipped, err)：skipped 为被 L1 拒绝/跳过的条数（观测用）。
func (s *OracleServiceV2) collectSettlementItems() ([]settlementItem, int, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, 0, err
	}

	items := make([]settlementItem, 0, len(users))
	skipped := 0
	for _, user := range users {
		userService, err := s.userServiceRepo.GetActiveByUserID(user.ID)
		if err != nil || userService == nil || !userService.IsActive {
			continue
		}

		// 用量（模拟）→ 计价（确定可复算，替换原 rand.Intn 随机金额，design §4.2）。
		// 单位锁：dataMB=MB、callMin=minute。
		dataMB, callMin, err := s.operatorAPI.GetUsage(user.ID, userService.OperatorID)
		if err != nil {
			skipped++
			continue
		}
		// L1 计价 + usage 上界 + 单 bill 硬上限（design §7.0①）。超界/超额：拒绝该条 + 告警，不出账。
		amount6, err := s.pricing.Price(dataMB, callMin, userService.OperatorID)
		if err != nil {
			log.Printf("WARN: 计价拒绝(L1)：user=%d op=%d dataMB=%d callMin=%d err=%v（跳过该 bill）", user.ID, userService.OperatorID, dataMB, callMin, err)
			skipped++
			continue
		}
		// amount6 == 0（低于最小出账额）→ 跳过（合约侧 amount>0 才 createBill，design §4.1）。
		if amount6.Sign() == 0 {
			continue
		}

		items = append(items, settlementItem{
			userID: user.ID,
			SettlementItem: SettlementItem{
				User:       common.HexToAddress(user.WalletAddr),
				OperatorID: userService.OperatorID,
				Amount6:    amount6,
			},
		})
	}
	return items, skipped, nil
}

// SettleMonthlyOnChain 是 T6 月度结算上链主入口（design §6.1）：收集计价项 → 组批 ≤25 →
// L2 闸（绝对闸 + 月均熔断 + 冷启动回退）→ 逐批 client.MonthlySettlement → month+batchIndex 幂等。
// 幂等可重跑（已 confirmed 批不重发，失败批续跑）。未注入 settlement 编排器（无 owner key/未配链）→ error。
func (s *OracleServiceV2) SettleMonthlyOnChain(ctx context.Context, month string) (*SettlementSummary, error) {
	if s.settlement == nil {
		return nil, fmt.Errorf("结算编排器未装配（owner key 未注入或链未配置）：月度上链结算不可用")
	}
	collected, _, err := s.collectSettlementItems()
	if err != nil {
		return nil, err
	}
	items := make([]SettlementItem, len(collected))
	for i, c := range collected {
		items[i] = c.SettlementItem
	}
	return s.settlement.SettleMonth(ctx, month, items)
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