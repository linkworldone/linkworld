package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
)

// 事件 topic0（keccak256(签名)）。
//
// design §5.2 / §6.3 铁律：字段解码一律走 abigen `*Filterer.Parse*`（类型安全），本文件仅供
// 轻量过滤 / 日志用，不参与字段解码。topic0 取自冻结 ABI 注释的权威值（与 bindings 中 abigen
// 注释的 event ID 一致），避免手写签名字符串再 keccak 出错（历史 bug：BillCreated 写成 5 参）。
//
// 与 internal/blockchain/bindings 的 abigen 注释逐一核对（FilterXxx 函数上方 "binding the
// contract event 0x..." 即该事件 topic0）。
var (
	// UserRegistered(address indexed user, string email, uint256 tokenId)
	UserRegisteredTopic = common.HexToHash("0x89105a1c6a3c2fbd471255c66a31ccab604af5697f67d7e2e9a0028c5e4dbd91")

	// DepositMade(address indexed user, uint256 amount)
	DepositMadeTopic = common.HexToHash("0xd15c9547ea5c06670c0010ce19bc32d54682a4b3801ece7f3ab0c3f17106b4bb")

	// DepositWithdrawn(address indexed user, uint256 principal, uint256 interest)
	DepositWithdrawnTopic = common.HexToHash("0x7719804546c0185709e60c90d164447ff251a5ba29af0216faa921350f6bebf7")

	// BillCreated(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 platformFee)
	// 注意：4 参（旧 signatures.go 错写 5 参，topic 匹配不到）。第三参 totalAmount = amount + platformFee 含费总额。
	BillCreatedTopic = common.HexToHash("0xcdfdeecd9f301cb609cbfd87c3a7f1e4d3da395ff0ba5084da583d3b6deced21")

	// BillPaid(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 operatorAmount)
	BillPaidTopic = common.HexToHash("0x53646f88205e1dd1de6fdaa8898d840fbe17129cb98729533d3c83a3a0e6045a")

	// TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)（Deposit 合约 emit）
	TrafficCardMintedTopic = common.HexToHash("0x3072026fc6418657755f94a3ef0972a28bb4266b72de75d96c3e999d4ff7067d")

	// CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)（TrafficCardNFT 合约 emit）
	CardMintedTopic = common.HexToHash("0xcb1c56d5745b05695241c17b7cfaece9a64f70bc32508c0997990633be8056ff")

	// SimRedeemed(address indexed user, uint256 daysCount, uint256[] tokenIds)
	// 流量卡销毁兑换 SIM（新玩法）：daysCount=销毁卡数=兑换天数；user indexed，daysCount/tokenIds 在 data 区。
	// 替代已移除的 ServiceActivated。单张 burn 也 emit SimRedeemed(...,1,...)。
	SimRedeemedTopic = common.HexToHash("0x2b22b8bfbca140907fa0a889499309fc6098870b0feb3b4af1993c079007b6b7")

	// TrafficCardApplied(uint256 indexed billId)（桩事件，仅记录不改金额）
	TrafficCardAppliedTopic = common.HexToHash("0x6ee1062e7611525ff44a2bea1e6ccffff047903a965c69aa96334ad680cd8701")

	// ESimRedeemed(address indexed user, uint256 tokenId, string activationCode, string smDpAddress)
	ESimRedeemedTopic = common.HexToHash("0x05ee3a1aa9b5412da9d4c0680a87e47a6014024f963aeaa94bc6d831d4d03cf4")

	// UsageDataSubmitted(address indexed user, uint256 operatorId, uint256 amount)
	// 只有 user indexed；operatorId/amount 在 data 区（design §6.3 解码歧义澄清）。
	UsageDataSubmittedTopic = common.HexToHash("0x02511b012a361a308c6f48148bcde643071e85e39257f525372ff1161aaf7c9b")
)
