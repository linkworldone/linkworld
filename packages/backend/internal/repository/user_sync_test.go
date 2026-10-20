package repository

import (
	"strings"
	"testing"
	"time"

	"linkworld-backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSyncFromChainRefreshesExistingUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := NewUserRepository(db)
	wallet := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	if err := repo.Create(&models.User{WalletAddr: wallet, IsActive: true}); err != nil {
		t.Fatalf("create API profile: %v", err)
	}
	registeredAt := time.Unix(1_800_000_000, 0).UTC()
	if err := repo.SyncFromChain(&models.User{
		WalletAddr:   strings.ToUpper(wallet),
		Email:        "chain@example.test",
		TokenID:      42,
		IsActive:     true,
		RegisteredAt: registeredAt,
	}); err != nil {
		t.Fatalf("sync chain registration: %v", err)
	}
	got, err := repo.FindByWallet(wallet)
	if err != nil {
		t.Fatalf("find synced user: %v", err)
	}
	if got.TokenID != 42 || got.Email != "chain@example.test" || !got.RegisteredAt.Equal(registeredAt) {
		t.Fatalf("chain fields not synced: token=%d email=%q registeredAt=%v", got.TokenID, got.Email, got.RegisteredAt)
	}
}
