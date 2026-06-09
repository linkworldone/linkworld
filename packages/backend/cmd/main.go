package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"linkworld-backend/internal/blockchain"
	"linkworld-backend/internal/config"
	"linkworld-backend/internal/handlers"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
	"linkworld-backend/internal/services"
	"linkworld-backend/internal/sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Operator{},
		&models.Deposit{},
		&models.Bill{},
		&models.UserService{},
		&models.UsageData{},
	)

	// Seed operators if empty
	var count int64
	db.Model(&models.Operator{}).Count(&count)
	if count == 0 {
		operators := []models.Operator{
			{Name: "T-Mobile", Region: "United States", CountryCode: "US", RequiredDeposit: "0.01", IsActive: true},
			{Name: "Vodafone", Region: "United Kingdom", CountryCode: "GB", RequiredDeposit: "0.008", IsActive: true},
			{Name: "Orange", Region: "France", CountryCode: "FR", RequiredDeposit: "0.008", IsActive: true},
			{Name: "MTS", Region: "Russia", CountryCode: "RU", RequiredDeposit: "0.005", IsActive: true},
			{Name: "SoftBank", Region: "Japan", CountryCode: "JP", RequiredDeposit: "0.012", IsActive: true},
			{Name: "Viettel", Region: "Vietnam", CountryCode: "VN", RequiredDeposit: "0.003", IsActive: true},
			{Name: "Unitel", Region: "Laos", CountryCode: "LA", RequiredDeposit: "0.003", IsActive: true},
			{Name: "Smart", Region: "Cambodia", CountryCode: "KH", RequiredDeposit: "0.003", IsActive: true},
			{Name: "AIS", Region: "Thailand", CountryCode: "TH", RequiredDeposit: "0.004", IsActive: true},
			{Name: "Maxis", Region: "Malaysia", CountryCode: "MY", RequiredDeposit: "0.004", IsActive: true},
			{Name: "Globe", Region: "Philippines", CountryCode: "PH", RequiredDeposit: "0.003", IsActive: true},
		}
		for _, op := range operators {
			db.Create(&op)
		}
		log.Println("Seeded 11 operators")
	}

	userRepo := repository.NewUserRepository(db)
	operatorRepo := repository.NewOperatorRepository(db)
	billRepo := repository.NewBillRepository(db)
	userServiceRepo := repository.NewUserServiceRepository(db)
	depositRepo := repository.NewDepositRepository(db)
	usageRepo := repository.NewUsageDataRepository(db)

	userService := services.NewUserService(userRepo)
	operatorService := services.NewOperatorService(operatorRepo)
	billingService := services.NewBillingService(billRepo, userRepo)
	oracleService := services.NewOracleService(userServiceRepo, billRepo, usageRepo, userRepo)
	notificationService := services.NewNotificationService(userRepo)
	depositService := services.NewDepositService(depositRepo, userRepo)
	userServiceService := services.NewUserServiceService(userServiceRepo, operatorRepo, userRepo)
	virtualGen := services.NewVirtualNumberGenerator()
	operatorAPI := services.NewOperatorAPISimulator()
	oracleV2 := services.NewOracleServiceV2(operatorAPI, userRepo, userServiceRepo, billRepo, usageRepo)
	usageService := services.NewUsageService(oracleV2, usageRepo, userRepo, userServiceRepo)

	// Initialize blockchain sync (optional - runs in background)
	if deployments, err := config.LoadDeployments("configs/deployments.json"); err == nil {
		// RPC 单一优先级（design §6.5）：deployments.json.rpcUrl 为准，env RPC_URL 覆盖。
		rpcURL := deployments.ResolveRPCURL(os.Getenv("RPC_URL"))
		// chainID 一致校验（design §6.5 / arch-review）：deployments.chainId 必须与 env CHAIN_ID 一致。
		if envChainID := parseUint64(os.Getenv("CHAIN_ID")); envChainID != 0 {
			if cerr := deployments.ValidateChainID(envChainID); cerr != nil {
				log.Fatal(cerr)
			}
		}
		// 占位零地址合约提示（T4 event_sync 据此跳过订阅）。
		if ph := deployments.PlaceholderContracts(); len(ph) > 0 {
			log.Printf("WARN: 合约占位零地址（未上链，事件同步将跳过）：%v", ph)
		}
		if rpcURL != "" {
			bcClient, err := blockchain.NewClient(rpcURL, deployments.ChainID, deployments.Proxies)
			if err == nil {
				// owner key 内存注入开启链上写（design §6.2 / §7.1）：从 env 读，仅内存、不落盘、
				// 不进日志（client 内仅打 owner address）；缺失或 chainID 不一致 → 写降级关闭，
				// 只读 + event_sync 仍可跑（owner=平台 root，不 fatal）。
				if ownerKey := os.Getenv("ORACLE_OWNER_PRIVATE_KEY"); ownerKey != "" {
					if werr := bcClient.EnableOwnerWrites(ownerKey); werr != nil {
						log.Printf("WARN: 链上写降级关闭（EnableOwnerWrites 失败）：%v", werr)
					}
				} else {
					log.Println("WARN: 未设置 ORACLE_OWNER_PRIVATE_KEY，链上写降级关闭（只读仍可）")
				}
				if ethClient := bcClient.EthClient(); ethClient != nil {
					eventSync := sync.NewEventSync(ethClient, userRepo, deployments.Proxies)
					go eventSync.Start(context.Background())
				}
				log.Println("Blockchain event sync started")
			}
		}
	}

	handler := handlers.NewHandler(userService, operatorService, billingService, oracleService, notificationService, depositService, userServiceService, virtualGen, oracleV2, usageService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/api/register", handler.Register)
	r.GET("/api/user/:wallet", handler.GetUser)
	r.GET("/api/operators", handler.GetOperators)
	r.POST("/api/service/activate", handler.ActivateService)
	r.POST("/api/service/deactivate", handler.DeactivateService)
	r.GET("/api/service/:wallet", handler.GetUserService)
	r.GET("/api/bills/:wallet", handler.GetBills)
	r.POST("/api/bills/pay", handler.PayBill)
	r.POST("/api/deposit", handler.Deposit)
	r.GET("/api/deposit/:wallet", handler.GetDeposit)
	r.GET("/api/deposit/:wallet/history", handler.GetDepositHistory)
	r.POST("/api/withdraw", handler.Withdraw)
	r.POST("/api/virtual-number/generate", handler.GenerateVirtualNumber)
	r.GET("/api/countries", handler.GetCountryList)
	r.GET("/api/usage/:wallet", handler.GetUsage)
	r.POST("/api/oracle/monthly-bill", handler.TriggerMonthlyBill)
	r.POST("/api/usage/submit", handler.SubmitUsage)
	r.POST("/api/notification/send", handler.SendNotification)

	r.Run(":8080")
}

// parseUint64 解析 env 字符串为 uint64；空或非法返回 0（视为「未设置」，由调用方决定是否强制）。
func parseUint64(s string) uint64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
