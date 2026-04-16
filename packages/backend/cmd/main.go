package main

import (
	"log"
	"math/rand"
	"time"

	"linkworld-backend/internal/config"
	"linkworld-backend/internal/handlers"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
	"linkworld-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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
	oracleService := services.NewOracleService(userServiceRepo, billRepo)
	notificationService := services.NewNotificationService(userRepo)
	depositService := services.NewDepositService(depositRepo, userRepo)
	userServiceService := services.NewUserServiceService(userServiceRepo, operatorRepo, userRepo)
	virtualGen := services.NewVirtualNumberGenerator()
	operatorAPI := services.NewOperatorAPISimulator()
	oracleV2 := services.NewOracleServiceV2(operatorAPI, userRepo, billRepo, usageRepo)
	usageService := services.NewUsageService(oracleV2, usageRepo, userRepo)

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
	r.POST("/api/withdraw", handler.Withdraw)
	r.POST("/api/virtual-number/generate", handler.GenerateVirtualNumber)
	r.GET("/api/countries", handler.GetCountryList)
	r.GET("/api/usage/:wallet", handler.GetUsage)
	r.POST("/api/oracle/monthly-bill", handler.TriggerMonthlyBill)

	r.Run(":8080")
}
