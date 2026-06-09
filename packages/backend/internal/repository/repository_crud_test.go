package repository

// T8 repository 关键路径 CRUD（design §8「repository CRUD 关键路径」+ B2/B3 对账记账契约）。
//
// 聚焦资损相关的两条记账契约（其余 CRUD 已由上层 service/handler/sync 测试间接覆盖）：
//   REPO-01：GetTotalByUserID 只计 confirmed deposit —— pending/withdraw 均不计入余额（B3 单一路径）
//   REPO-02：SetPayIntent 只写意向不动终态 —— 命中写 PayIntentTxHash；未命中 bill 返回 ErrRecordNotFound（B2）

import (
	"strings"
	"testing"

	"linkworld-backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func crudDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Deposit{}, &models.Bill{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// REPO-01：余额 = 仅 confirmed deposit；pending deposit 与所有 withdraw 不计入。
func TestRepo_GetTotalByUserID_ConfirmedDepositOnly(t *testing.T) {
	db := crudDB(t)
	r := NewDepositRepository(db)
	const uid = uint(1)

	// confirmed deposit 1_000_000（计入）
	db.Create(&models.Deposit{UserID: uid, Amount: "1000000", Type: "deposit", Status: models.DepositStatusConfirmed})
	// pending deposit 9_000_000（不计入——尚未 K 块确认，B5）
	db.Create(&models.Deposit{UserID: uid, Amount: "9000000", Type: "deposit", Status: models.DepositStatusPending})
	// pending withdraw 500_000（不计入——withdraw 终态唯一由事件回填，B3）
	db.Create(&models.Deposit{UserID: uid, Amount: "500000", Type: "withdraw", Status: models.DepositStatusPending})
	// 另一用户的 confirmed deposit（隔离，不应串到 uid）
	db.Create(&models.Deposit{UserID: uint(2), Amount: "7777777", Type: "deposit", Status: models.DepositStatusConfirmed})

	total, err := r.GetTotalByUserID(uid)
	if err != nil {
		t.Fatalf("GetTotalByUserID: %v", err)
	}
	if total != "1000000" {
		t.Fatalf("REPO-01 🔴: 余额只应计 confirmed deposit=1000000（pending/withdraw/他人均排除），got %s", total)
	}
}

// REPO-01b：无任何记录时余额为 "0"（COALESCE 兜底，不空串/不报错）。
func TestRepo_GetTotalByUserID_EmptyIsZero(t *testing.T) {
	db := crudDB(t)
	r := NewDepositRepository(db)
	total, err := r.GetTotalByUserID(uint(99))
	if err != nil {
		t.Fatalf("GetTotalByUserID: %v", err)
	}
	if total != "0" {
		t.Fatalf("REPO-01b: 无记录余额应为 0，got %q", total)
	}
}

// REPO-02：SetPayIntent 命中 → 写 PayIntentTxHash 不动 IsPaid；未命中 → ErrRecordNotFound。
func TestRepo_SetPayIntent_IntentOnly(t *testing.T) {
	db := crudDB(t)
	r := NewBillRepository(db)
	const uid = uint(1)

	bill := &models.Bill{UserID: uid, OperatorID: 1, Amount: "1000000", IsPaid: false}
	db.Create(bill)

	if err := r.SetPayIntent(uid, bill.ID, "0xintent"); err != nil {
		t.Fatalf("REPO-02: 命中应成功，err=%v", err)
	}
	var got models.Bill
	db.First(&got, bill.ID)
	if got.PayIntentTxHash != "0xintent" {
		t.Fatalf("REPO-02: 应写 PayIntentTxHash，got %q", got.PayIntentTxHash)
	}
	if got.IsPaid {
		t.Fatal("REPO-02 🔴: SetPayIntent 绝不能置 IsPaid（终态唯一由 event_sync 回填）")
	}

	// 未命中（bill 不存在或不属于该 user）→ ErrRecordNotFound（防越权改他人 bill）。
	if err := r.SetPayIntent(uid, 99999, "0xnope"); err != gorm.ErrRecordNotFound {
		t.Fatalf("REPO-02: 未命中应返回 ErrRecordNotFound，got %v", err)
	}
	if err := r.SetPayIntent(uint(2), bill.ID, "0xnotmine"); err != gorm.ErrRecordNotFound {
		t.Fatalf("REPO-02 🔴: 他人 user 改本 bill 应未命中（user_id 绑定），got %v", err)
	}
}
