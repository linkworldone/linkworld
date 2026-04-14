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
