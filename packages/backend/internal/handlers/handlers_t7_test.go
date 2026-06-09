package handlers

// T7 端点加固 TDD（先红后绿）：bills/pay 降级、withdraw 降级、usage/submit 范围校验。
//
// 覆盖用例：
//   PAY-INTENT-01：bills/pay → 写 PayIntentTxHash，IsPaid 仍 false（终态唯一由 event_sync 回填）
//   WD-01：withdraw 不凭 txHash 记账（写 pending 意向，不计入余额 GetTotalByUserID）
//   USAGE-01：usage/submit 超 max → 400（gin binding max= 范围校验，B4）

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"linkworld-backend/internal/middleware"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
	"linkworld-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func t7DB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Bill{}, &models.Deposit{}, &models.UsageData{},
		&models.Operator{}, &models.UserService{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// t7Handler 构造仅装配 T7 相关 service 的 handler（其余传 nil，T7 用例不触达）。
func t7Handler(db *gorm.DB) *Handler {
	userRepo := repository.NewUserRepository(db)
	billRepo := repository.NewBillRepository(db)
	depositRepo := repository.NewDepositRepository(db)
	usageRepo := repository.NewUsageDataRepository(db)
	userServiceRepo := repository.NewUserServiceRepository(db)

	billingService := services.NewBillingService(billRepo, userRepo)
	depositService := services.NewDepositService(depositRepo, userRepo)
	oracleService := services.NewOracleService(userServiceRepo, billRepo, usageRepo, userRepo)

	return &Handler{
		billingService: billingService,
		depositService: depositService,
		oracleService:  oracleService,
	}
}

func init() { gin.SetMode(gin.TestMode) }

func postJSON(r http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// stubWalletAuth 模拟 WalletAuth 中间件：把已验证钱包写入 CtxWallet（裸挂 handler 的 T7 用例用，
// 不跑真签名链路）。R1 修复后 handler 操作主体唯一取 CtxWallet，故裸挂测试必须显式注入。
func stubWalletAuth(wallet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxWallet, wallet)
		c.Next()
	}
}

// --- PAY-INTENT-01 ---

func TestPayBill_WritesIntent_NotPaid(t *testing.T) {
	db := t7DB(t)
	h := t7Handler(db)

	user := &models.User{WalletAddr: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Email: "a@b.c", IsActive: true}
	db.Create(user)
	bill := &models.Bill{UserID: user.ID, OperatorID: 1, Amount: "1000000", IsPaid: false}
	db.Create(bill)

	r := gin.New()
	r.POST("/api/bills/pay", stubWalletAuth(user.WalletAddr), h.PayBill)

	w := postJSON(r, "/api/bills/pay", gin.H{
		"bill_id": bill.ID,
		"tx_hash": "0xdeadbeef",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("PAY-INTENT-01: 应 200，got %d body=%s", w.Code, w.Body.String())
	}

	var got models.Bill
	db.First(&got, bill.ID)
	if got.IsPaid {
		t.Fatal("PAY-INTENT-01 🔴: bills/pay 绝不能置 IsPaid=true（终态唯一由 event_sync 回填）")
	}
	if got.PayIntentTxHash != "0xdeadbeef" {
		t.Fatalf("PAY-INTENT-01: 应写 PayIntentTxHash 意向，got %q", got.PayIntentTxHash)
	}
}

// --- WD-01 ---

func TestWithdraw_PendingIntent_NotCountedInBalance(t *testing.T) {
	db := t7DB(t)
	h := t7Handler(db)

	user := &models.User{WalletAddr: "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", Email: "b@b.c", IsActive: true}
	db.Create(user)
	// 已确认押金 5000000，余额应为 5000000。
	db.Create(&models.Deposit{UserID: user.ID, Amount: "5000000", Type: "deposit", Status: models.DepositStatusConfirmed})

	r := gin.New()
	r.POST("/api/withdraw", stubWalletAuth(user.WalletAddr), h.Withdraw)

	w := postJSON(r, "/api/withdraw", gin.H{
		"tx_hash": "0xfeed",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("WD-01: 应 200，got %d body=%s", w.Code, w.Body.String())
	}

	// 提现 pending 意向不得计入余额（GetTotalByUserID 只认 confirmed deposit）。
	depRepo := repository.NewDepositRepository(db)
	total, err := depRepo.GetTotalByUserID(user.ID)
	if err != nil {
		t.Fatalf("GetTotalByUserID: %v", err)
	}
	if total != "5000000" {
		t.Fatalf("WD-01: 提现 pending 意向不得改余额，期望 5000000，got %s", total)
	}

	// 若落了 withdraw 记录，必须是 pending（不能是 confirmed 凭前端 txHash 直接记账）。
	var withdraws []models.Deposit
	db.Where("user_id = ? AND type = ?", user.ID, "withdraw").Find(&withdraws)
	for _, wd := range withdraws {
		if wd.Status == models.DepositStatusConfirmed {
			t.Fatal("WD-01 🔴: withdraw 不得凭前端 txHash 写 confirmed 记账（终态唯一由 DepositWithdrawn 事件回填）")
		}
	}
}

// --- USAGE-01 ---

func TestSubmitUsage_OverMax_400(t *testing.T) {
	db := t7DB(t)
	h := t7Handler(db)

	user := &models.User{WalletAddr: "0xcccccccccccccccccccccccccccccccccccccccc", Email: "c@b.c", IsActive: true}
	db.Create(user)

	r := gin.New()
	r.POST("/api/usage/submit", h.SubmitUsage)

	// data_usage 超出 gin binding max= 上界 → 400。
	w := postJSON(r, "/api/usage/submit", gin.H{
		"wallet":      user.WalletAddr,
		"operator_id": 1,
		"data_usage":  uint64(1) << 62, // 天文用量，必超上界
		"call_usage":  100,
		"signature":   "sig",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("USAGE-01: 超 max 应 400，got %d body=%s", w.Code, w.Body.String())
	}
}

func TestSubmitUsage_WithinMax_OK(t *testing.T) {
	db := t7DB(t)
	h := t7Handler(db)

	user := &models.User{WalletAddr: "0xdddddddddddddddddddddddddddddddddddddddd", Email: "d@b.c", IsActive: true}
	db.Create(user)

	r := gin.New()
	r.POST("/api/usage/submit", h.SubmitUsage)

	w := postJSON(r, "/api/usage/submit", gin.H{
		"wallet":      user.WalletAddr,
		"operator_id": 1,
		"data_usage":  10000,
		"call_usage":  500,
		"signature":   "sig",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("USAGE-01: 合法范围应 200，got %d body=%s", w.Code, w.Body.String())
	}
}

var _ = fmt.Sprintf
