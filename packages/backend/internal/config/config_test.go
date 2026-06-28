package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// writeTempJSON 把内容写到临时文件，返回路径（测试结束自动清理）。
func writeTempJSON(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "deployments.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp json: %v", err)
	}
	return p
}

// fullDeployment 是一份带 7 proxies + usdt + usdtDecimals + abiHash 的合法 421614 schema。
const fullDeployment = `{
  "chainId": 421614,
  "rpcUrl": "https://sepolia-rollup.arbitrum.io/rpc",
  "proxies": {
    "FeeManager": "0xF9d4777b760cc3a0F39eE0E11Cc936E34dcfc033",
    "UserRegistry": "0x0D0E7AeB3437682964d8164835eAE31c86451268",
    "ServiceManager": "0x82CB050c84F3BBEfC01D089d8579805Eb493BA14",
    "TrafficCardNFT": "0x8B29aC425eD0b021CFFb308494707A5f4e6DEd31",
    "Payment": "0x85Ffe2f47dF883982A6c98f665670e045fd0bfd9",
    "Deposit": "0x1c73baEceE72d0867b046f939Dd27fbbc714332b",
    "Oracle": "0x1820f818dF0dE96d29eA3AA7007785eBE46662D1"
  },
  "usdt": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "usdtDecimals": 6,
  "abiHash": {
    "Oracle": "0xd1d0f83aaff6bb08fb11bbcefa961d55ceef7b6d6f8fe86b272772b20c25f359"
  }
}`

// TestLoadDeployments_ParsesProxiesKey 验证键名 bug 已修复：JSON 键 `proxies` 能被
// 反序列化进 Proxies（旧代码 struct tag 是 proxies 但 JSON 写 contracts → 永远空 map）。
func TestLoadDeployments_ParsesProxiesKey(t *testing.T) {
	path := writeTempJSON(t, fullDeployment)
	d, err := LoadDeployments(path)
	if err != nil {
		t.Fatalf("LoadDeployments: %v", err)
	}

	if len(d.Proxies) != 7 {
		t.Fatalf("Proxies should have 7 entries, got %d (键名 bug 未修复?)", len(d.Proxies))
	}
	want := common.HexToAddress("0x1820f818dF0dE96d29eA3AA7007785eBE46662D1")
	if got := d.Proxies["Oracle"]; got != want {
		t.Errorf("Proxies[Oracle] = %s, want %s", got.Hex(), want.Hex())
	}
}

// TestLoadDeployments_ParsesExtendedFields 验证 struct 扩展字段 Usdt/UsdtDecimals/AbiHash
// 正确反序列化 + ChainID/RpcURL。
func TestLoadDeployments_ParsesExtendedFields(t *testing.T) {
	path := writeTempJSON(t, fullDeployment)
	d, err := LoadDeployments(path)
	if err != nil {
		t.Fatalf("LoadDeployments: %v", err)
	}

	if d.ChainID != 421614 {
		t.Errorf("ChainID = %d, want 421614", d.ChainID)
	}
	if d.RpcURL != "https://sepolia-rollup.arbitrum.io/rpc" {
		t.Errorf("RpcURL = %q", d.RpcURL)
	}
	wantUsdt := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	if d.Usdt != wantUsdt {
		t.Errorf("Usdt = %s, want %s", d.Usdt.Hex(), wantUsdt.Hex())
	}
	if d.UsdtDecimals != 6 {
		t.Errorf("UsdtDecimals = %d, want 6", d.UsdtDecimals)
	}
	if d.AbiHash["Oracle"] != "0xd1d0f83aaff6bb08fb11bbcefa961d55ceef7b6d6f8fe86b272772b20c25f359" {
		t.Errorf("AbiHash[Oracle] = %q", d.AbiHash["Oracle"])
	}
}

// TestIsPlaceholder 验证占位零地址判断：零地址/空字符串视为占位，非零真实地址不是。
func TestIsPlaceholder(t *testing.T) {
	cases := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"explicit zero hex", common.HexToAddress("0x0000000000000000000000000000000000000000"), true},
		{"real address", common.HexToAddress("0x1820f818dF0dE96d29eA3AA7007785eBE46662D1"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPlaceholder(tc.addr); got != tc.want {
				t.Errorf("IsPlaceholder(%s) = %v, want %v", tc.addr.Hex(), got, tc.want)
			}
		})
	}
}

// TestDeployments_PlaceholderContracts 验证占位地址识别能力（供 T4 event_sync 跳过零地址合约）。
func TestDeployments_PlaceholderContracts(t *testing.T) {
	partial := `{
  "chainId": 421614,
  "rpcUrl": "https://sepolia-rollup.arbitrum.io/rpc",
  "proxies": {
    "FeeManager": "0x0000000000000000000000000000000000000000",
    "Oracle": "0x1820f818dF0dE96d29eA3AA7007785eBE46662D1"
  },
  "usdt": "0x0000000000000000000000000000000000000000",
  "usdtDecimals": 6
}`
	path := writeTempJSON(t, partial)
	d, err := LoadDeployments(path)
	if err != nil {
		t.Fatalf("LoadDeployments: %v", err)
	}

	if !d.HasPlaceholders() {
		t.Error("HasPlaceholders() = false, want true (FeeManager + usdt 为零地址)")
	}
	placeholders := d.PlaceholderContracts()
	if len(placeholders) != 1 || placeholders[0] != "FeeManager" {
		t.Errorf("PlaceholderContracts() = %v, want [FeeManager]", placeholders)
	}
	if IsPlaceholder(d.Proxies["Oracle"]) {
		t.Error("Oracle should not be placeholder")
	}
	if !IsPlaceholder(d.Usdt) {
		t.Error("Usdt should be placeholder")
	}
}

// TestValidateChainID 验证 chainID 一致校验（design §6.5 / arch-review）：
// deployments chainId 与 env CHAIN_ID 不一致须报错（防连 A 链发 B 链）。
func TestValidateChainID(t *testing.T) {
	d := &Deployments{ChainID: 421614}

	if err := d.ValidateChainID(421614); err != nil {
		t.Errorf("ValidateChainID(421614) should pass, got %v", err)
	}
	if err := d.ValidateChainID(31337); err == nil {
		t.Error("ValidateChainID(31337) should fail (chainId 不一致), got nil")
	}
	// envChainID == 0 表示未设置 → 不强制校验（以 deployments 为准）。
	if err := d.ValidateChainID(0); err != nil {
		t.Errorf("ValidateChainID(0) should pass (env 未设), got %v", err)
	}
}

// TestResolveRPCURL 验证 RPC 来源单一优先级（design §6.5）：json.rpcUrl 为准，env RPC_URL 覆盖。
func TestResolveRPCURL(t *testing.T) {
	d := &Deployments{RpcURL: "https://json-rpc.example"}

	// env 为空 → 用 json。
	if got := d.ResolveRPCURL(""); got != "https://json-rpc.example" {
		t.Errorf("ResolveRPCURL(\"\") = %q, want json url", got)
	}
	// env 非空 → 覆盖 json。
	if got := d.ResolveRPCURL("https://env-rpc.example"); got != "https://env-rpc.example" {
		t.Errorf("ResolveRPCURL(env) = %q, want env url", got)
	}
}

// TestRealDeploymentsJSON 校验仓库内真实 configs/deployments.json：
// chainId=421614、7 个 proxies 为真实地址、usdt 为真实地址、usdtDecimals=6。
func TestRealDeploymentsJSON(t *testing.T) {
	d, err := LoadDeployments("../../configs/deployments.json")
	if err != nil {
		t.Fatalf("LoadDeployments(real): %v", err)
	}

	if d.ChainID != 421614 {
		t.Errorf("ChainID = %d, want 421614 (0G 16602 残留?)", d.ChainID)
	}
	if d.RpcURL == "https://evm-testnet.0g.ai" {
		t.Error("rpcUrl 仍是 0G 残留")
	}
	wantContracts := []string{"MockUSDT", "FeeManager", "UserRegistry", "ServiceManager", "TrafficCardNFT", "Payment", "Deposit", "Oracle"}
	if len(d.Proxies) != len(wantContracts) {
		t.Fatalf("Proxies len = %d, want %d", len(d.Proxies), len(wantContracts))
	}
	for _, name := range wantContracts {
		addr, ok := d.Proxies[name]
		if !ok {
			t.Errorf("Proxies 缺 %s", name)
			continue
		}
		if IsPlaceholder(addr) {
			t.Errorf("Proxies[%s] = %s，本轮应为真实地址", name, addr.Hex())
		}
	}
	if IsPlaceholder(d.Usdt) {
		t.Errorf("usdt = %s，本轮应为真实地址", d.Usdt.Hex())
	}
	if d.UsdtDecimals != 6 {
		t.Errorf("UsdtDecimals = %d, want 6", d.UsdtDecimals)
	}
}
