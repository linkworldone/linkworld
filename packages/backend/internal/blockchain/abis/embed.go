// Package abis 把 8 份合约的原始 ABI JSON（从 packages/contracts/artifacts 提取）
// 嵌入二进制，供 abihash 校验与绑定加载离线使用。
//
// 重新生成：cd packages/backend && bash scripts/gen-bindings.sh
// （脚本同时刷新 bindings/ 与本目录的 *.json）。
package abis

import "embed"

//go:embed *.json
var FS embed.FS

// Names 是本目录嵌入的 8 个合约 ABI 文件基名（不含 .json），顺序与 deployments
// abiHash 一致；MockUSDT 无 abiHash（部署侧未纳入指纹）。
var Names = []string{
	"FeeManager",
	"UserRegistry",
	"ServiceManager",
	"TrafficCardNFT",
	"Payment",
	"Deposit",
	"Oracle",
	"MockUSDT",
}
