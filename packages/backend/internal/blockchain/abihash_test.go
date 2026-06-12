package blockchain

import (
	"testing"

	"linkworld-backend/internal/blockchain/abis"
)

// expectedABIHash 来自合约侧 packages/contracts/deployments/hardhat.json 的 abiHash 字段，
// 由 deploy.ts 用 keccak256(utf8(ethers.Interface.formatJson())) 生成，是后端比对的权威基准。
// MockUSDT 不在 abiHash 内（部署侧未纳入指纹），故不列入。
var expectedABIHash = map[string]string{
	"FeeManager":     "0x75813eaf919f9b6e95ee77044d9c26e3027695cceefb31c9f05d17fd7df948a1",
	"UserRegistry":   "0xca527da9916960e91d34757251357a77f061befb5f3e778e7cdbc748b79675d8",
	"ServiceManager": "0x8a634a3b4d0e3e5c0b015500b9bcb4d26499f52ccafd06c6c03a687b0259e370",
	"TrafficCardNFT": "0xf74ea0faa6f4a2c72d1a6efd9b2858620fdba1200cf5947dd0d826d5411558fd",
	"Payment":        "0xa8579606c85a8c91da901f8400127767da1e93557daf056e11d123b438118356",
	"Deposit":        "0xe876a431188d76a43cb1600aced605d2dca7f68d7516bb053dfc6362b23d37d9",
	"Oracle":         "0xd1d0f83aaff6bb08fb11bbcefa961d55ceef7b6d6f8fe86b272772b20c25f359",
}

// TestABIHashMatchesDeployments 校验嵌入的 7 个业务合约 ABI 用本地 formatJson 复刻算法
// 算出的 abiHash 与 hardhat.json 完全一致——证明后端绑定与链上部署 ABI 同步。
func TestABIHashMatchesDeployments(t *testing.T) {
	for name, want := range expectedABIHash {
		raw, err := abis.FS.ReadFile(name + ".json")
		if err != nil {
			t.Fatalf("read embedded abi %s: %v", name, err)
		}
		got, err := ComputeABIHash(raw)
		if err != nil {
			t.Fatalf("compute abihash %s: %v", name, err)
		}
		if got != want {
			t.Errorf("abiHash mismatch for %s:\n  got  %s\n  want %s", name, got, want)
		}
		if err := VerifyABIHash(name, raw, want); err != nil {
			t.Errorf("VerifyABIHash(%s): %v", name, err)
		}
	}
}

// TestVerifyABIHashDetectsMismatch 确认指纹不一致时 VerifyABIHash 报错（资损/不同步防线）。
func TestVerifyABIHashDetectsMismatch(t *testing.T) {
	raw, err := abis.FS.ReadFile("Oracle.json")
	if err != nil {
		t.Fatalf("read Oracle abi: %v", err)
	}
	bogus := "0x0000000000000000000000000000000000000000000000000000000000000000"
	if err := VerifyABIHash("Oracle", raw, bogus); err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
}

// TestMockUSDTHashStable 仅确认 MockUSDT ABI 能被规范化+哈希（无回归即可，无部署侧基准比对）。
func TestMockUSDTHashStable(t *testing.T) {
	raw, err := abis.FS.ReadFile("MockUSDT.json")
	if err != nil {
		t.Fatalf("read MockUSDT abi: %v", err)
	}
	if _, err := ComputeABIHash(raw); err != nil {
		t.Fatalf("compute MockUSDT abihash: %v", err)
	}
}
