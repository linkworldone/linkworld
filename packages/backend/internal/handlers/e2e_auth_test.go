package handlers

// T8 端到端鉴权流（design §8 + arch-review §七 + task-07 遗留）。
//
// 与 handlers_t7_test.go（裸挂 handler 测降级语义）和 middleware_test.go（裸挂中间件测签名）
// 的区别：本文件把 middleware.NewWalletAuth + 受保护 handler 串成 HTTP 全链路——
// 用 middleware.WalletAuthDigest 构造真实 EIP-712 签名头，经中间件 ecrecover + 消费 nonce 后
// 才到达 handler，验证「签名→鉴权→业务→降级落库」整条链路在受保护端点上成立。
//
// 覆盖用例：
//   E2E-AUTH-01：带有效签名 → PayBill 全链路 200 + 写 PayIntentTxHash 不置 IsPaid（B2）
//   E2E-AUTH-02：缺签名头 → 401，handler 不执行（Bill 不被改）
//   E2E-AUTH-03：nonce 重放 → 二次 401（一次性台账消费，跨中间件+handler 链路有效）
//   E2E-AUTH-04：带有效签名 → Withdraw 全链路 200 + pending 不计入余额（B3）

import (
	"crypto/ecdsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"linkworld-backend/internal/middleware"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const e2eChainID uint64 = 421614

func e2eDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Bill{}, &models.Deposit{}, &models.UsageData{},
		&models.Operator{}, &models.UserService{}, &models.WalletNonce{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func e2eKey(t *testing.T) (*ecdsa.PrivateKey, string) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return key, crypto.PubkeyToAddress(key.PublicKey).Hex()
}

// e2eSign 用私钥对 (wallet,nonce,action,chainId) 的 EIP-712 摘要签名（前端等价行为）。
func e2eSign(t *testing.T, key *ecdsa.PrivateKey, wallet, nonce, action string) string {
	t.Helper()
	digest, err := middleware.WalletAuthDigest(wallet, nonce, action, e2eChainID)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	sig, err := crypto.Sign(digest, key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	const hexchars = "0123456789abcdef"
	out := make([]byte, len(sig)*2)
	for i, v := range sig {
		out[i*2] = hexchars[v>>4]
		out[i*2+1] = hexchars[v&0x0f]
	}
	return "0x" + string(out)
}

// e2eSignedReq 构造一个带 WalletAuth 签名头的 POST 请求（全链路驱动）。
func e2eSignedReq(method, path, body, wallet, nonce, action, sig string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(middleware.HeaderWalletAddr, wallet)
	req.Header.Set(middleware.HeaderWalletNonce, nonce)
	req.Header.Set(middleware.HeaderWalletAction, action)
	req.Header.Set(middleware.HeaderWalletSig, sig)
	return req
}

// e2eRouter 装配「WalletAuth 中间件 + 受保护端点」的真实路由（与 main.go §6.6 挂载方式一致）。
func e2eRouter(db *gorm.DB) (*gin.Engine, *repository.WalletNonceRepository) {
	h := t7Handler(db)
	nonceRepo := repository.NewWalletNonceRepository(db)
	wa := middleware.NewWalletAuth(nonceRepo, e2eChainID)

	r := gin.New()
	r.POST("/api/bills/pay", wa, h.PayBill)
	r.POST("/api/withdraw", wa, h.Withdraw)
	return r, nonceRepo
}

// --- E2E-AUTH-01：有效签名穿过中间件到达 PayBill，写意向不置 IsPaid ---

func TestE2E_WalletAuth_PayBill_FullChain(t *testing.T) {
	db := e2eDB(t)
	r, nonceRepo := e2eRouter(db)
	key, addr := e2eKey(t)

	db.Create(&models.User{WalletAddr: addr, Email: "e@b.c", IsActive: true})
	var u models.User
	db.First(&u, "wallet_addr = ?", addr)
	bill := &models.Bill{UserID: u.ID, OperatorID: 1, Amount: "1000000", IsPaid: false}
	db.Create(bill)

	nonce, err := nonceRepo.Issue(addr)
	if err != nil {
		t.Fatalf("issue nonce: %v", err)
	}
	sig := e2eSign(t, key, addr, nonce, "bills/pay")

	body := `{"wallet":"` + addr + `","bill_id":` + itoa(bill.ID) + `,"tx_hash":"0xpayintent"}`
	w := httptest.NewRecorder()
	r.ServeHTTP(w, e2eSignedReq(http.MethodPost, "/api/bills/pay", body, addr, nonce, "bills/pay", sig))

	if w.Code != http.StatusOK {
		t.Fatalf("E2E-AUTH-01: 有效签名全链路应 200，got %d body=%s", w.Code, w.Body.String())
	}
	var got models.Bill
	db.First(&got, bill.ID)
	if got.IsPaid {
		t.Fatal("E2E-AUTH-01 🔴: 全链路下 bills/pay 仍绝不能置 IsPaid（B2 终态唯一由 event_sync 回填）")
	}
	if got.PayIntentTxHash != "0xpayintent" {
		t.Fatalf("E2E-AUTH-01: 全链路应写 PayIntentTxHash 意向，got %q", got.PayIntentTxHash)
	}
}

// --- E2E-AUTH-02：缺签名头 → 401，handler 不执行 ---

func TestE2E_WalletAuth_MissingSig_BlocksHandler(t *testing.T) {
	db := e2eDB(t)
	r, _ := e2eRouter(db)

	db.Create(&models.User{WalletAddr: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Email: "a@b.c", IsActive: true})
	var u models.User
	db.First(&u, "wallet_addr = ?", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	bill := &models.Bill{UserID: u.ID, OperatorID: 1, Amount: "1000000", IsPaid: false}
	db.Create(bill)

	// 无任何 WalletAuth 头，直接打受保护端点。
	body := `{"wallet":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","bill_id":` + itoa(bill.ID) + `,"tx_hash":"0xshouldnotwrite"}`
	req := httptest.NewRequest(http.MethodPost, "/api/bills/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("E2E-AUTH-02: 缺签名头应 401，got %d", w.Code)
	}
	var got models.Bill
	db.First(&got, bill.ID)
	if got.PayIntentTxHash != "" {
		t.Fatalf("E2E-AUTH-02 🔴: 鉴权失败时 handler 不应执行（PayIntentTxHash 不应被写），got %q", got.PayIntentTxHash)
	}
}

// --- E2E-AUTH-03：nonce 重放在全链路下二次 401 ---

func TestE2E_WalletAuth_ReplayNonce_FullChain(t *testing.T) {
	db := e2eDB(t)
	r, nonceRepo := e2eRouter(db)
	key, addr := e2eKey(t)

	db.Create(&models.User{WalletAddr: addr, Email: "r@b.c", IsActive: true})
	var u models.User
	db.First(&u, "wallet_addr = ?", addr)
	bill := &models.Bill{UserID: u.ID, OperatorID: 1, Amount: "1000000", IsPaid: false}
	db.Create(bill)

	nonce, _ := nonceRepo.Issue(addr)
	sig := e2eSign(t, key, addr, nonce, "bills/pay")
	body := `{"wallet":"` + addr + `","bill_id":` + itoa(bill.ID) + `,"tx_hash":"0xonce"}`

	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, e2eSignedReq(http.MethodPost, "/api/bills/pay", body, addr, nonce, "bills/pay", sig))
	if w1.Code != http.StatusOK {
		t.Fatalf("E2E-AUTH-03 前置: 首次应 200，got %d body=%s", w1.Code, w1.Body.String())
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, e2eSignedReq(http.MethodPost, "/api/bills/pay", body, addr, nonce, "bills/pay", sig))
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("E2E-AUTH-03 🔴: 同 nonce 重放在全链路下二次应 401（一次性台账），got %d", w2.Code)
	}
}

// --- E2E-AUTH-04：有效签名穿中间件到 Withdraw，pending 不计入余额 ---

func TestE2E_WalletAuth_Withdraw_FullChain(t *testing.T) {
	db := e2eDB(t)
	r, nonceRepo := e2eRouter(db)
	key, addr := e2eKey(t)

	db.Create(&models.User{WalletAddr: addr, Email: "w@b.c", IsActive: true})
	var u models.User
	db.First(&u, "wallet_addr = ?", addr)
	db.Create(&models.Deposit{UserID: u.ID, Amount: "5000000", Type: "deposit", Status: models.DepositStatusConfirmed})

	nonce, _ := nonceRepo.Issue(addr)
	sig := e2eSign(t, key, addr, nonce, "withdraw")
	body := `{"wallet":"` + addr + `","tx_hash":"0xwd"}`

	w := httptest.NewRecorder()
	r.ServeHTTP(w, e2eSignedReq(http.MethodPost, "/api/withdraw", body, addr, nonce, "withdraw", sig))
	if w.Code != http.StatusOK {
		t.Fatalf("E2E-AUTH-04: 有效签名全链路应 200，got %d body=%s", w.Code, w.Body.String())
	}

	depRepo := repository.NewDepositRepository(db)
	total, err := depRepo.GetTotalByUserID(u.ID)
	if err != nil {
		t.Fatalf("GetTotalByUserID: %v", err)
	}
	if total != "5000000" {
		t.Fatalf("E2E-AUTH-04 🔴: 全链路下提现 pending 意向不得改余额，期望 5000000，got %s", total)
	}
}

// itoa 是最小整数转字符串（避免引 strconv 仅为拼 JSON）。
func itoa(u uint) string {
	if u == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	return string(buf[i:])
}
