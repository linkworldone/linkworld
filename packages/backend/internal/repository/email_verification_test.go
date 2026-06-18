package repository

// 邮箱绑定验证码台账测试（与 WalletNonce 同一次性/防爆破/限流安全模式）。
// 覆盖：6 位数字码、过期拒绝、attempts≥5 锁定、一次性消费、限流。
// 用 sqlite in-memory（glebarez/sqlite，与 repository_crud_test.go 同策略）。

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"linkworld-backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func emailDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.EmailVerification{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

const (
	testWallet = "0xAbC0000000000000000000000000000000000001"
	testEmail  = "User@Example.com"
)

var sixDigits = regexp.MustCompile(`^\d{6}$`)

// EV-01：Issue 产出 6 位纯数字码，落库 wallet 小写归一化、未用、未过期。
func TestEmailVerification_Issue_SixDigitCode(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	code, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if !sixDigits.MatchString(code) {
		t.Fatalf("EV-01: 验证码应为 6 位纯数字，got %q", code)
	}

	var rec models.EmailVerification
	if err := db.Where("code = ?", code).First(&rec).Error; err != nil {
		t.Fatalf("EV-01: 验证码应落库，err=%v", err)
	}
	if rec.Wallet != strings.ToLower(testWallet) {
		t.Fatalf("EV-01: wallet 应小写归一化，got %q", rec.Wallet)
	}
	if rec.Used {
		t.Fatal("EV-01: 新签发码应 Used=false")
	}
	if !rec.ExpiresAt.After(time.Now()) {
		t.Fatal("EV-01: 新签发码应未过期")
	}
}

// EV-02：Verify 一次性 —— 第一次正确 true，第二次同码 false（已消费 Used=true）。
func TestEmailVerification_Verify_OneTimeUse(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	code, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if !r.Verify(testWallet, testEmail, code) {
		t.Fatal("EV-02: 首次正确验证应 true")
	}
	if r.Verify(testWallet, testEmail, code) {
		t.Fatal("EV-02 🔴: 同一码二次验证必须 false（一次性消费防重放）")
	}
}

// EV-02b：Verify 大小写不敏感匹配 email/wallet（与 Issue 归一化口径一致）。
func TestEmailVerification_Verify_CaseInsensitive(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	code, _ := r.Issue(strings.ToLower(testWallet), strings.ToLower(testEmail))
	// 用大写 wallet/email 来验证应仍命中（归一化）。
	if !r.Verify(strings.ToUpper(testWallet), strings.ToUpper(testEmail), code) {
		t.Fatal("EV-02b: wallet/email 大小写应归一化匹配")
	}
}

// EV-03：过期码拒绝 —— 手动把 ExpiresAt 改到过去，Verify 必 false。
func TestEmailVerification_Verify_ExpiredRejected(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	code, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	// 强制过期。
	if err := db.Model(&models.EmailVerification{}).
		Where("code = ?", code).
		Update("expires_at", time.Now().Add(-time.Minute)).Error; err != nil {
		t.Fatalf("force expire: %v", err)
	}
	if r.Verify(testWallet, testEmail, code) {
		t.Fatal("EV-03 🔴: 过期验证码必须拒绝（false）")
	}
}

// EV-04：attempts≥5 锁定 —— 连续 5 次错码后，即便给出正确码也拒绝（记录已锁死）。
func TestEmailVerification_Verify_LockedAfterMaxAttempts(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	code, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	// 5 次错码（每次 attempts++）。
	for i := 0; i < 5; i++ {
		if r.Verify(testWallet, testEmail, "000000") {
			t.Fatalf("EV-04: 错码第 %d 次不应通过", i+1)
		}
	}
	var rec models.EmailVerification
	db.Where("code = ?", code).First(&rec)
	if rec.Attempts < 5 {
		t.Fatalf("EV-04: 5 次错码后 attempts 应 ≥5，got %d", rec.Attempts)
	}
	// 此后即使给正确码也应拒绝（attempts<5 条件不再满足，锁死）。
	if r.Verify(testWallet, testEmail, code) {
		t.Fatal("EV-04 🔴: attempts≥5 锁死后即便正确码也必须拒绝（防爆破）")
	}
}

// EV-05：限流 —— 同 wallet 60 秒内重复 Issue 返回 ErrEmailRateLimited。
func TestEmailVerification_Issue_RateLimited(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	if _, err := r.Issue(testWallet, testEmail); err != nil {
		t.Fatalf("首次 Issue: %v", err)
	}
	_, err := r.Issue(testWallet, testEmail)
	if err != ErrEmailRateLimited {
		t.Fatalf("EV-05 🔴: 60 秒内重复签发应限流（ErrEmailRateLimited），got %v", err)
	}
}

// EV-05b：限流窗口外可重发 —— 把上一条 CreatedAt 拨到 61 秒前，新 Issue 应成功并作废旧码。
func TestEmailVerification_Issue_ResendAfterInterval(t *testing.T) {
	db := emailDB(t)
	r := NewEmailVerificationRepository(db)

	oldCode, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("首次 Issue: %v", err)
	}
	// 把旧记录 CreatedAt 拨到限流窗口外。
	if err := db.Model(&models.EmailVerification{}).
		Where("code = ?", oldCode).
		Update("created_at", time.Now().Add(-61*time.Second)).Error; err != nil {
		t.Fatalf("backdate: %v", err)
	}

	newCode, err := r.Issue(testWallet, testEmail)
	if err != nil {
		t.Fatalf("EV-05b: 窗口外重发应成功，err=%v", err)
	}
	if newCode == oldCode {
		t.Fatal("EV-05b: 新码应不同于旧码")
	}
	// 旧码应已被作废（Used=true），不能再验证通过。
	if r.Verify(testWallet, testEmail, oldCode) {
		t.Fatal("EV-05b 🔴: 重发后旧码应被作废，不可再用")
	}
	// 新码可正常验证。
	if !r.Verify(testWallet, testEmail, newCode) {
		t.Fatal("EV-05b: 新码应可验证通过")
	}
}
