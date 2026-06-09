package middleware

// T7 鉴权中间件 TDD（先红后绿）。
//
// 覆盖用例：
//   ADMIN-01：正确 key 过 / 错误 key 401 / 缺 key（NewAdminAuth 返回 error，启动 fail）
//   WALLET-01：有效 EIP-712 签名 + 未用 nonce → 过
//   WALLET-02（🔴 重放）：同 nonce 二次用 → 拒绝（一次性台账消费式）
//   WALLET-03：错 chainId/domain → 拒绝
//   WALLET-04：签名 wallet != 请求 wallet → 拒绝
//
// 链交互取舍（与 T3/T4 一致）：用内存 sqlite（纯 Go，无 CGO）承载 nonce 台账，
// EIP-712 签名用 crypto.Sign 对 TypedData hash 直接签（不依赖部署合约）。

import (
	"crypto/ecdsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testChainID uint64 = 421614

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.WalletNonce{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newTestKey(t *testing.T) (*ecdsa.PrivateKey, string) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
	return key, addr
}

func init() {
	gin.SetMode(gin.TestMode)
}

// --- ADMIN-01 ---

func TestAdminAuth_MissingKey_StartupFail(t *testing.T) {
	if _, err := NewAdminAuth(""); err == nil {
		t.Fatal("ADMIN-01: 缺 ADMIN_API_KEY 时 NewAdminAuth 必须返回 error（启动 fail），不静默放行")
	}
}

func TestAdminAuth_CorrectKey_Passes(t *testing.T) {
	mw, err := NewAdminAuth("secret-key")
	if err != nil {
		t.Fatalf("NewAdminAuth: %v", err)
	}
	r := gin.New()
	r.POST("/x", mw, func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Header.Set(HeaderAdminKey, "secret-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("ADMIN-01: 正确 key 应过，got %d", w.Code)
	}
}

func TestAdminAuth_WrongKey_401(t *testing.T) {
	mw, _ := NewAdminAuth("secret-key")
	r := gin.New()
	r.POST("/x", mw, func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Header.Set(HeaderAdminKey, "wrong-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ADMIN-01: 错误 key 应 401，got %d", w.Code)
	}
}

func TestAdminAuth_NoHeader_401(t *testing.T) {
	mw, _ := NewAdminAuth("secret-key")
	r := gin.New()
	r.POST("/x", mw, func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ADMIN-01: 缺 header 应 401，got %d", w.Code)
	}
}

// --- WalletAuth 辅助 ---

// signWallet 用私钥对 (wallet, nonce, action) 的 EIP-712 TypedData hash 签名，返回 0x 前缀十六进制签名。
func signWallet(t *testing.T, key *ecdsa.PrivateKey, wallet string, nonce string, action string, chainID uint64) string {
	t.Helper()
	digest, err := WalletAuthDigest(wallet, nonce, action, chainID)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	sig, err := crypto.Sign(digest, key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return "0x" + toHex(sig)
}

func toHex(b []byte) string {
	const hexchars = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexchars[v>>4]
		out[i*2+1] = hexchars[v&0x0f]
	}
	return string(out)
}

func issueNonce(t *testing.T, repo *repository.WalletNonceRepository, wallet string) string {
	t.Helper()
	n, err := repo.Issue(wallet)
	if err != nil {
		t.Fatalf("issue nonce: %v", err)
	}
	return n
}

func walletRouter(t *testing.T, db *gorm.DB, chainID uint64) (*gin.Engine, *repository.WalletNonceRepository) {
	return walletRouterAction(t, db, chainID, "bills/pay")
}

// walletRouterAction 装配绑定指定 expectedAction 的 WalletAuth 路由（action 绑端点测试用）。
func walletRouterAction(t *testing.T, db *gorm.DB, chainID uint64, action string) (*gin.Engine, *repository.WalletNonceRepository) {
	repo := repository.NewWalletNonceRepository(db)
	mw := NewWalletAuth(repo, chainID, action)
	r := gin.New()
	r.POST("/x", mw, func(c *gin.Context) {
		// 中间件应把已验证 wallet 放进 context，供 handler 取用。
		w, _ := c.Get(CtxWallet)
		c.JSON(http.StatusOK, gin.H{"wallet": w})
	})
	return r, repo
}

func walletReq(wallet, nonce, action, sig string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderWalletAddr, wallet)
	req.Header.Set(HeaderWalletNonce, nonce)
	req.Header.Set(HeaderWalletAction, action)
	req.Header.Set(HeaderWalletSig, sig)
	return req
}

// --- WALLET-01：有效签名 + 未用 nonce → 过 ---

func TestWalletAuth_ValidSignature_Passes(t *testing.T) {
	db := newTestDB(t)
	r, repo := walletRouter(t, db, testChainID)
	key, addr := newTestKey(t)

	nonce := issueNonce(t, repo, addr)
	sig := signWallet(t, key, addr, nonce, "bills/pay", testChainID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addr, nonce, "bills/pay", sig))
	if w.Code != http.StatusOK {
		t.Fatalf("WALLET-01: 有效签名 + 未用 nonce 应过，got %d body=%s", w.Code, w.Body.String())
	}
}

// --- WALLET-02（🔴 重放）：同 nonce 二次用 → 拒绝 ---

func TestWalletAuth_ReplayNonce_Rejected(t *testing.T) {
	db := newTestDB(t)
	r, repo := walletRouter(t, db, testChainID)
	key, addr := newTestKey(t)

	nonce := issueNonce(t, repo, addr)
	sig := signWallet(t, key, addr, nonce, "bills/pay", testChainID)

	// 首次：通过（消费 nonce）。
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, walletReq(addr, nonce, "bills/pay", sig))
	if w1.Code != http.StatusOK {
		t.Fatalf("WALLET-02 前置: 首次应过，got %d", w1.Code)
	}

	// 二次：同 nonce 同签名重放 → 必须拒绝（一次性台账已作废）。
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, walletReq(addr, nonce, "bills/pay", sig))
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("WALLET-02 🔴 重放: 同 nonce 二次用必须 401（一次性消费台账），got %d", w2.Code)
	}
}

// --- WALLET-03：错 chainId → 拒绝 ---

func TestWalletAuth_WrongChainID_Rejected(t *testing.T) {
	db := newTestDB(t)
	// 服务端绑定 testChainID，但签名用错误 chainID → digest 不匹配 → ecrecover 出别的地址 → 拒绝。
	r, repo := walletRouter(t, db, testChainID)
	key, addr := newTestKey(t)

	nonce := issueNonce(t, repo, addr)
	sig := signWallet(t, key, addr, nonce, "bills/pay", testChainID+999)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addr, nonce, "bills/pay", sig))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("WALLET-03: 错 chainId 应 401，got %d", w.Code)
	}
}

// --- WALLET-04：签名 wallet != 请求 wallet → 拒绝 ---

func TestWalletAuth_WalletMismatch_Rejected(t *testing.T) {
	db := newTestDB(t)
	r, repo := walletRouter(t, db, testChainID)
	keyA, _ := newTestKey(t)
	_, addrB := newTestKey(t)

	// 用 A 的私钥签名，但请求声称是 B 的 wallet。
	nonce := issueNonce(t, repo, addrB)
	sig := signWallet(t, keyA, addrB, nonce, "bills/pay", testChainID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addrB, nonce, "bills/pay", sig))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("WALLET-04: 签名地址 != 请求 wallet 应 401，got %d", w.Code)
	}
}

// --- 未签发 nonce → 拒绝（台账无此 nonce） ---

func TestWalletAuth_UnknownNonce_Rejected(t *testing.T) {
	db := newTestDB(t)
	r, _ := walletRouter(t, db, testChainID)
	key, addr := newTestKey(t)

	sig := signWallet(t, key, addr, "never-issued-nonce", "bills/pay", testChainID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addr, "never-issued-nonce", "bills/pay", sig))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("未签发 nonce 应 401，got %d", w.Code)
	}
}

// --- WALLET-ACTION-01：为 action="withdraw" 签的 nonce 用于 bills/pay 端点 → 拒绝 ---
//
// review Medium：action 进 EIP-712 摘要但未校验匹配端点。绑定后，端点 expectedAction="bills/pay"
// 而签名/头里 action="withdraw" → action mismatch，401，且不消费 nonce。
func TestWalletAuth_ActionMismatch_Rejected(t *testing.T) {
	db := newTestDB(t)
	// 端点绑定 bills/pay。
	r, repo := walletRouterAction(t, db, testChainID, "bills/pay")
	key, addr := newTestKey(t)

	// 但用户为 withdraw 动作签名（合法签名、有效 nonce），挪用到 bills/pay 端点。
	nonce := issueNonce(t, repo, addr)
	sig := signWallet(t, key, addr, nonce, "withdraw", testChainID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addr, nonce, "withdraw", sig))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("WALLET-ACTION-01: 为 withdraw 签的 nonce 用于 bills/pay 端点必须 401（action 绑端点），got %d", w.Code)
	}
	// action 不匹配应在消费 nonce 前拒绝：nonce 仍可被正确动作使用（未被作废）。
	sigOK := signWallet(t, key, addr, nonce, "bills/pay", testChainID)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, walletReq(addr, nonce, "bills/pay", sigOK))
	if w2.Code != http.StatusOK {
		t.Fatalf("WALLET-ACTION-01: action 不匹配不应消费 nonce，正确 action 重试应过，got %d body=%s", w2.Code, w2.Body.String())
	}
}

// --- WALLET-ACTION：空 action 头 → 拒绝（防绕过绑定）---
func TestWalletAuth_EmptyAction_Rejected(t *testing.T) {
	db := newTestDB(t)
	r, repo := walletRouterAction(t, db, testChainID, "bills/pay")
	key, addr := newTestKey(t)

	nonce := issueNonce(t, repo, addr)
	sig := signWallet(t, key, addr, nonce, "", testChainID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, walletReq(addr, nonce, "", sig))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("空 action 头应 401（action 绑端点），got %d", w.Code)
	}
}
