package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"linkworld-backend/internal/blockchain"
	"linkworld-backend/internal/blockchain/bindings"
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

	// Seed operators if empty —— 显式写 ID=链上 operatorId(1..11)（design §4.5 / arch-review）。
	// operatorId 固定映射：models.Operator.ID == 链上 operatorId，不靠 name 比对。
	// SeedOperators() 是 seed 与 PricingService 费率表的单一事实源（避免两处漂移）。
	var count int64
	db.Model(&models.Operator{}).Count(&count)
	if count == 0 {
		for _, op := range services.SeedOperators() {
			op := op // 显式主键，逐个 Create（GORM 允许指定 ID）
			db.Create(&op)
		}
		log.Println("Seeded 11 operators（ID=链上 operatorId 1..11，固定映射）")
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
	// 计价引擎（design §4.2）：费率表占位 + usage 上界 L1 + 单 bill 硬上限 L1（构造/init 已 log.Warn 占位刺眼提示）。
	pricingService := services.NewPricingService()
	oracleV2 := services.NewOracleServiceV2(operatorAPI, pricingService, userRepo, userServiceRepo, billRepo, usageRepo)
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
					// operatorId 固定映射 sanity check（design §4.5 / arch-review）：读链 ServiceManager
					// 校验 seed 的 operatorId 在链上存在、countryCode 一致、分账地址非零，不靠 name 比对。
					// 占位零地址（未上链）→ 降级跳过（不阻塞启动）。不一致 → 默认 warn（资金环境可配 fail-fast）。
					smAddr := deployments.Proxies["ServiceManager"]
					if config.IsPlaceholder(smAddr) {
						log.Println("WARN: ServiceManager 占位零地址（未上链），跳过 operatorId sanity check")
					} else if smCaller, cerr := bindings.NewServiceManagerCaller(smAddr, ethClient); cerr != nil {
						log.Printf("WARN: 构造 ServiceManager caller 失败，跳过 operatorId sanity check：%v", cerr)
					} else {
						reader := services.NewServiceManagerCallerReader(smCaller)
						if mism, scerr := services.SanityCheckOperators(reader, services.SeedOperators()); scerr != nil {
							log.Printf("WARN: operatorId sanity check 读链失败（降级跳过）：%v", scerr)
						} else if len(mism) > 0 {
							// 默认 warn（design §4.5：可配置降级为 fail-fast，资金环境建议 fail）。
							log.Printf("WARN: operatorId 链上/seed 不一致（共 %d 项，分账打错风险，上线前必查）：%v", len(mism), mism)
						} else {
							log.Println("operatorId sanity check 通过（seed ID==链上 operatorId，countryCode 一致，分账地址非零）")
						}
					}

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
