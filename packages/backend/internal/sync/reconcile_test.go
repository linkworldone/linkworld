package sync

// T8 对账金额不一致告警路径（design §6.3/§7.5 检测层 + arch-review B 对账）。
//
// processUsageDataSubmitted 是「检测层」对账：链上 UsageDataSubmitted(user,operatorId,amount)
// 与后端计价 Bill.Amount 比对，不一致 → WARN 告警（资损信号）。codegraph 确认该 process 无覆盖测试。
// 本文件补：① 金额不一致 → 走告警路径；② 金额一致 → 不告警；③ 非资金事件即时 confirmed（不等 K 块）。
//
// 复用 event_sync_test.go 的 fakeSource/buildLog/newSyncWithFake/seedUser 基础设施。

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"strings"
	"testing"

	"linkworld-backend/internal/blockchain/bindings"
	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// captureLog 捕获 std log 输出供断言告警路径（log 包是 processUsageDataSubmitted 的告警出口）。
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	orig := log.Writer()
	flags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(orig)
		log.SetFlags(flags)
	}()
	fn()
	return buf.String()
}

// usageLog 构造一条 UsageDataSubmitted（仅 user indexed；operatorId/amount 在 data 区）。
func usageLog(t *testing.T, blk uint64, idx uint, opID, amount *big.Int) types.Log {
	return buildLog(t, bindings.OracleMetaData.ABI, "UsageDataSubmitted", testOracle, blk, idx,
		[]common.Hash{common.BytesToHash(testWallet.Bytes())}, opID, amount)
}

func TestReconcile_AmountMismatch_Alerts(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	// 后端计价出账 1_000_000；链上 UsageDataSubmitted 入账 999_999（不一致 1）。
	billRepo := repository.NewBillRepository(db)
	billRepo.Create(&models.Bill{UserID: uid, OperatorID: 7, Amount: "1000000", IsPaid: false})

	src.head = 5
	src.logsByBlk[5] = []types.Log{usageLog(t, 5, 0, big.NewInt(7), big.NewInt(999_999))}

	out := captureLog(t, func() {
		if err := s.SyncOnce(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	if !strings.Contains(out, "对账金额不一致") {
		t.Fatalf("RECONCILE-01 🔴: 链上金额 != 后端计价应走告警路径（WARN 对账金额不一致），日志=%q", out)
	}
	if !strings.Contains(out, "后端=1000000") || !strings.Contains(out, "链上=999999") {
		t.Fatalf("RECONCILE-01: 告警应含双方金额（后端/链上），日志=%q", out)
	}
}

func TestReconcile_AmountMatch_NoAlert(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	uid := seedUser(t, ur)

	// 后端计价 == 链上入账（500_000），不应告警。
	billRepo := repository.NewBillRepository(db)
	billRepo.Create(&models.Bill{UserID: uid, OperatorID: 7, Amount: "500000", IsPaid: false})

	src.head = 5
	src.logsByBlk[5] = []types.Log{usageLog(t, 5, 0, big.NewInt(7), big.NewInt(500_000))}

	out := captureLog(t, func() {
		if err := s.SyncOnce(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
	if strings.Contains(out, "对账金额不一致") {
		t.Fatalf("RECONCILE-02 🔴: 金额一致不应告警，误报日志=%q", out)
	}
}

// RECONCILE-03：UsageDataSubmitted 是非资金终态事件，应即时 confirmed（不进 seen 等 K 块）。
func TestReconcile_UsageDataSubmitted_ImmediateConfirmed(t *testing.T) {
	db := newTestDB(t)
	src := newFakeSource()
	s, ur := newSyncWithFake(t, src, db)
	seedUser(t, ur)

	src.head = 5 // 深度远不足 K(3)，但非资金事件不受两阶段约束。
	src.logsByBlk[5] = []types.Log{usageLog(t, 5, 0, big.NewInt(7), big.NewInt(500_000))}
	if err := s.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	var ev models.ChainEvent
	if err := db.First(&ev, "event_name = ?", "UsageDataSubmitted").Error; err != nil {
		t.Fatalf("应落 UsageDataSubmitted ChainEvent: %v", err)
	}
	if ev.Status != models.ChainEventStatusConfirmed {
		t.Fatalf("RECONCILE-03: 非资金事件 UsageDataSubmitted 应即时 confirmed，得到 %s", ev.Status)
	}
}
