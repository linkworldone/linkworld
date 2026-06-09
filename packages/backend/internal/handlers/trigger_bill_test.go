package handlers

// T8 月度结算 handler 降级序列化（design §6.5 + task-06/07 遗留「SettlementSummary=null 降级分支」）。
//
// TriggerMonthlyBill 两段：① 始终落 DB Bill；② 结算编排器已装配才上链。
// 未装配（无 owner key/未配链）→ 返回 settlement: null 降级（不报错，HTTP 200）。
// 本测覆盖 codegraph 标记「无覆盖」的 TriggerMonthlyBill + SettlementSummary 序列化分支。
//
// 覆盖用例：
//   TRIGGER-01：未装配编排器 → 200 + settlement 字段为 null（降级序列化，不 500）
//   TRIGGER-02：装配编排器（mock client）→ 200 + settlement 为非 null 摘要对象

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"
	"linkworld-backend/internal/services"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// triggerHandler 装配 TriggerMonthlyBill 所需的 usageService + oracleV2（settlement 可选注入）。
func triggerHandler(db *gorm.DB, settlement *services.SettlementOrchestrator) *Handler {
	userRepo := repository.NewUserRepository(db)
	billRepo := repository.NewBillRepository(db)
	usageRepo := repository.NewUsageDataRepository(db)
	userServiceRepo := repository.NewUserServiceRepository(db)

	oracleV2 := services.NewOracleServiceV2(
		&services.OperatorAPISimulator{}, services.NewPricingService(),
		userRepo, userServiceRepo, billRepo, usageRepo,
	)
	if settlement != nil {
		oracleV2.SetSettlementOrchestrator(settlement)
	}
	usageSvc := services.NewUsageService(oracleV2, usageRepo, userRepo, userServiceRepo)

	return &Handler{oracleV2: oracleV2, usageService: usageSvc}
}

// TRIGGER-01：未装配编排器 → settlement: null 降级（200，不 500）。
func TestTriggerMonthlyBill_NotConfigured_SettlementNull(t *testing.T) {
	db := e2eDB(t)
	h := triggerHandler(db, nil) // 不注入编排器

	r := gin.New()
	r.POST("/api/oracle/monthly-bill", h.TriggerMonthlyBill)

	w := postJSON(r, "/api/oracle/monthly-bill", gin.H{"month": "2026-06"})
	if w.Code != http.StatusOK {
		t.Fatalf("TRIGGER-01: 未装配编排器属预期降级，应 200 不 500，got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应: %v", err)
	}
	s, ok := resp["settlement"]
	if !ok {
		t.Fatal("TRIGGER-01: 响应应含 settlement 字段")
	}
	if strings.TrimSpace(string(s)) != "null" {
		t.Fatalf("TRIGGER-01 🔴: 未装配编排器时 settlement 必须序列化为 null，got %s", string(s))
	}
}

// TRIGGER-02：装配编排器（mock client）→ settlement 为非 null 摘要。
func TestTriggerMonthlyBill_Configured_SettlementSummary(t *testing.T) {
	db := e2eDB(t)

	// 装配一个总是成功的 mock 结算客户端 + 内存批次 store（经生产 repo store，覆盖序列化全路径）。
	repo := repository.NewSettlementBatchRepository(db)
	store, err := services.NewSettlementBatchRepoStore(repo)
	if err != nil {
		t.Fatalf("build store: %v", err)
	}
	orch := services.NewSettlementOrchestrator(okClient{}, store)
	h := triggerHandler(db, orch)

	// 无激活服务用户 → collectSettlementItems 空 → 0 批，但编排器已装配，settlement 应为非 null 摘要对象。
	r := gin.New()
	r.POST("/api/oracle/monthly-bill", h.TriggerMonthlyBill)

	w := postJSON(r, "/api/oracle/monthly-bill", gin.H{"month": "2026-06"})
	if w.Code != http.StatusOK {
		t.Fatalf("TRIGGER-02: 应 200，got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]json.RawMessage
	json.Unmarshal(w.Body.Bytes(), &resp)
	s := strings.TrimSpace(string(resp["settlement"]))
	if s == "null" || s == "" {
		t.Fatalf("TRIGGER-02 🔴: 装配编排器后 settlement 应为非 null 摘要对象，got %q", s)
	}
	// 摘要应可解出 Month 字段（序列化结构正确）。
	if !strings.Contains(s, "2026-06") {
		t.Fatalf("TRIGGER-02: settlement 摘要应含 month=2026-06，got %s", s)
	}
}

// okClient 是 services.SettlementClient 的最小实现（总是成功）。
type okClient struct{}

func (okClient) MonthlySettlement(ctx context.Context, users []common.Address, operatorIds, amounts []*big.Int) (*gethtypes.Receipt, error) {
	return &gethtypes.Receipt{Status: gethtypes.ReceiptStatusSuccessful, TxHash: common.BigToHash(big.NewInt(1))}, nil
}

var _ = models.SettlementBatch{}
