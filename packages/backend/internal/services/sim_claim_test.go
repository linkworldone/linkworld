package services

// SimService.Claim 的 eSIM / physical 交付方式分支测试。
// 内存 sqlite（纯 Go，无 CGO）搭最小 DB，先建 user 再 Claim，断言两个分支的字段。

import (
	"strings"
	"testing"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newSimTestService(t *testing.T) (*SimService, *models.User) {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Sim{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	user := &models.User{WalletAddr: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", IsActive: true}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	simRepo := repository.NewSimRepository(db)
	return NewSimService(simRepo, userRepo), user
}

func TestSimService_Claim_Esim(t *testing.T) {
	svc, user := newSimTestService(t)

	// esim 分支：不传 recipient/addressLine 也应成功，并生成激活链接。
	sim, err := svc.Claim(user.WalletAddr, "US", "", "", models.SimDeliveryEsim, 3, "0xtx")
	if err != nil {
		t.Fatalf("Claim esim: %v", err)
	}
	if sim.DeliveryType != models.SimDeliveryEsim {
		t.Errorf("DeliveryType = %q, want %q", sim.DeliveryType, models.SimDeliveryEsim)
	}
	if sim.ActivationURL == "" {
		t.Error("esim ActivationURL should be non-empty")
	}
	const prefix = "https://esim.linkworld.io/activate?token="
	if !strings.HasPrefix(sim.ActivationURL, prefix) {
		t.Errorf("ActivationURL = %q, want prefix %q", sim.ActivationURL, prefix)
	}
	// token 应为 16 字节 hex = 32 个十六进制字符。
	token := strings.TrimPrefix(sim.ActivationURL, prefix)
	if len(token) != 32 {
		t.Errorf("token length = %d, want 32 (16-byte hex)", len(token))
	}
	if sim.Status != models.SimStatusPending {
		t.Errorf("Status = %q, want %q", sim.Status, models.SimStatusPending)
	}
	if sim.Days != 3 {
		t.Errorf("Days = %d, want 3", sim.Days)
	}
}

func TestSimService_Claim_Physical(t *testing.T) {
	svc, user := newSimTestService(t)

	// physical 分支：保留 recipient/addressLine，ActivationURL 为空。
	sim, err := svc.Claim(user.WalletAddr, "JP", "Alice", "123 Main St", models.SimDeliveryPhysical, 2, "0xtx")
	if err != nil {
		t.Fatalf("Claim physical: %v", err)
	}
	if sim.DeliveryType != models.SimDeliveryPhysical {
		t.Errorf("DeliveryType = %q, want %q", sim.DeliveryType, models.SimDeliveryPhysical)
	}
	if sim.ActivationURL != "" {
		t.Errorf("physical ActivationURL = %q, want empty", sim.ActivationURL)
	}
	if sim.Recipient != "Alice" || sim.AddressLine != "123 Main St" {
		t.Errorf("recipient/address not preserved: %q / %q", sim.Recipient, sim.AddressLine)
	}
}

func TestSimService_Claim_DefaultsToPhysical(t *testing.T) {
	svc, user := newSimTestService(t)

	// 旧契约兼容：deliveryType 为空（旧前端不传）应按 physical 处理。
	sim, err := svc.Claim(user.WalletAddr, "US", "Bob", "456 Oak Ave", "", 1, "0xtx")
	if err != nil {
		t.Fatalf("Claim default: %v", err)
	}
	if sim.DeliveryType != models.SimDeliveryPhysical {
		t.Errorf("DeliveryType = %q, want %q", sim.DeliveryType, models.SimDeliveryPhysical)
	}
	if sim.ActivationURL != "" {
		t.Errorf("ActivationURL = %q, want empty", sim.ActivationURL)
	}
}
