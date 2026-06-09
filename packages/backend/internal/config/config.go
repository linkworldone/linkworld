package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Deployments 是 backend/configs/deployments.json 的反序列化结构，对齐合约侧
// packages/contracts/deployments/<net>.json 的范式（键名 `proxies`，含 usdt/usdtDecimals/abiHash）。
//
// 历史 bug（design §6.5 / arch-review）：旧 struct tag 读 `proxies` 但 JSON 写 `contracts`，
// 二者不一致 → Proxies 永远反序列化为空 map，event_sync 静默失效。T2 将 deployments.json
// 键名统一改为 `proxies`，与 hardhat.json 范式一致，消除双重不一致。
type Deployments struct {
	ChainID      uint64                    `json:"chainId"`
	RpcURL       string                    `json:"rpcUrl"`
	Proxies      map[string]common.Address `json:"proxies"`
	Usdt         common.Address            `json:"usdt"`
	UsdtDecimals uint8                     `json:"usdtDecimals"`
	AbiHash      map[string]string         `json:"abiHash"`
}

func LoadDeployments(path string) (*Deployments, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d Deployments
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if d.Proxies == nil {
		d.Proxies = map[string]common.Address{}
	}
	if d.AbiHash == nil {
		d.AbiHash = map[string]string{}
	}
	return &d, nil
}

// IsPlaceholder 判断地址是否为占位零地址。合约 1/3 真·上链前，421614 的合约地址
// 用零地址占位（待 PR#1 上链回填）。event_sync（T4）据此跳过零地址合约的事件过滤，
// 避免误匹配；client（T3）据此降级。
func IsPlaceholder(addr common.Address) bool {
	return addr == (common.Address{})
}

// HasPlaceholders 报告是否存在任一占位（零地址）合约或占位 usdt——即本套部署尚未完整上链。
func (d *Deployments) HasPlaceholders() bool {
	if IsPlaceholder(d.Usdt) {
		return true
	}
	for _, addr := range d.Proxies {
		if IsPlaceholder(addr) {
			return true
		}
	}
	return false
}

// PlaceholderContracts 返回 proxies 中所有占位（零地址）合约名，按入库顺序无保证；
// 供 T4 event_sync「零地址合约跳过订阅 + warn」与启动期诊断使用。
func (d *Deployments) PlaceholderContracts() []string {
	var names []string
	for name, addr := range d.Proxies {
		if IsPlaceholder(addr) {
			names = append(names, name)
		}
	}
	return names
}

// ValidateChainID 校验 deployments 的 chainId 与 env CHAIN_ID 一致（design §6.5 / arch-review）。
// 不一致 → 返回 error（防连 A 链发 B 链、签错链重放）。
// envChainID == 0 表示 env 未设置 → 以 deployments 为准，不强制校验。
func (d *Deployments) ValidateChainID(envChainID uint64) error {
	if envChainID == 0 {
		return nil
	}
	if d.ChainID != envChainID {
		return fmt.Errorf("chainId 不一致：deployments.json=%d, env CHAIN_ID=%d（连错网/签错链风险，拒绝启动）", d.ChainID, envChainID)
	}
	return nil
}

// ResolveRPCURL 统一 RPC 来源单一优先级（design §6.5）：deployments.json 的 rpcUrl 为准，
// env RPC_URL 非空时覆盖。明确优先级，避免连 A 链发 B 链。
func (d *Deployments) ResolveRPCURL(envRPCURL string) string {
	if envRPCURL != "" {
		return envRPCURL
	}
	return d.RpcURL
}

func InitDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "linkworld")
	dbname := getEnv("DB_NAME", "linkworld")

	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("PostgreSQL connected successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
