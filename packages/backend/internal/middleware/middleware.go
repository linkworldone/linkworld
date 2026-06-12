// Package middleware 提供 T7 敏感端点鉴权中间件（design §6.6 / arch-review 🔴 N1）。
//
// 两类鉴权（按操作主体划分）：
//   - AdminAuth：平台/管理触发端点（monthly-bill、notification/send、usage/submit）。
//     校验 X-Admin-Key 头 == ADMIN_API_KEY，常量时间比较（subtle.ConstantTimeCompare）。
//     缺 ADMIN_API_KEY → NewAdminAuth 返回 error，main.go 启动 fail（管理端点不允许裸奔）。
//   - WalletAuth：用户资金/写端点（bills/pay、withdraw、deposit、service 写）。
//     EIP-712 结构化签名 + ecrecover 还原签名地址 == 请求 wallet（绑定 msg.sender 语义）。
//     🔴 N1 防重放：服务端一次性 nonce 台账（消费式，签过即作废），绑定 chainId + domain，
//     禁纯 timestamp 时间窗。
package middleware

import (
	"crypto/subtle"
	"fmt"
	"math/big"
	"strings"

	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
)

// 请求头常量。
const (
	HeaderAdminKey     = "X-Admin-Key"
	HeaderWalletAddr   = "X-Wallet-Address"
	HeaderWalletNonce  = "X-Wallet-Nonce"
	HeaderWalletAction = "X-Wallet-Action"
	HeaderWalletSig    = "X-Wallet-Signature"
)

// CtxWallet 是 WalletAuth 校验通过后写入 gin.Context 的已验证钱包地址键（小写归一化）。
const CtxWallet = "auth_wallet"

// EIP-712 domain 常量（绑定 domain + chainId，arch-review 🔴 N1）。
const (
	eip712DomainName    = "LinkWorld"
	eip712DomainVersion = "1"
)

// NewAdminAuth 构造 AdminAuth 中间件。adminKey 为空 → 返回 error（main.go 据此启动 fail，
// 管理端点不允许裸奔，design §6.6）。
func NewAdminAuth(adminKey string) (gin.HandlerFunc, error) {
	if adminKey == "" {
		return nil, fmt.Errorf("ADMIN_API_KEY 未设置：管理端点鉴权不可用，拒绝启动（不允许裸奔）")
	}
	expected := []byte(adminKey)
	return func(c *gin.Context) {
		got := []byte(c.GetHeader(HeaderAdminKey))
		// 常量时间比较（防时序侧信道）。长度不等时 ConstantTimeCompare 返回 0。
		if subtle.ConstantTimeCompare(got, expected) != 1 {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: invalid admin key"})
			return
		}
		c.Next()
	}, nil
}

// WalletAuthDigest 计算 (wallet, nonce, action, chainId) 的 EIP-712 TypedData 摘要（待签 hash）。
// 前端用钱包私钥对该摘要签名（eth_signTypedData_v4 等价），后端据同一摘要 ecrecover。
// 绑定 chainId + domain，使跨链/跨域签名天然不可复用（WALLET-03）。
func WalletAuthDigest(wallet, nonce, action string, chainID uint64) ([]byte, error) {
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"WalletAuth": []apitypes.Type{
				{Name: "wallet", Type: "address"},
				{Name: "nonce", Type: "string"},
				{Name: "action", Type: "string"},
			},
		},
		PrimaryType: "WalletAuth",
		Domain: apitypes.TypedDataDomain{
			Name:    eip712DomainName,
			Version: eip712DomainVersion,
			ChainId: (*math.HexOrDecimal256)(new(big.Int).SetUint64(chainID)),
		},
		Message: apitypes.TypedDataMessage{
			"wallet": common.HexToAddress(wallet).Hex(),
			"nonce":  nonce,
			"action": action,
		},
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	// EIP-712 最终摘要：keccak256("\x19\x01" ‖ domainSeparator ‖ hashStruct(message))。
	raw := append([]byte{0x19, 0x01}, append(domainSeparator, messageHash...)...)
	return crypto.Keccak256(raw), nil
}

// NewWalletAuth 构造 WalletAuth 中间件。chainID 绑定到 EIP-712 domain（防跨链重放）。
// expectedAction 绑定本路由的语义动作（如 "bills/pay"/"withdraw"）：签名摘要里的 action 必须
// 严格等于 expectedAction（非空），否则拒绝——防止把为某动作签的 nonce 挪用到别的端点
// （WALLET-ACTION，review Medium：action 入摘要但未校验匹配端点）。expectedAction 为空属装配错误。
// 校验链路：取头部 (wallet, nonce, action, signature) → action 必须 == expectedAction
// → 重算 EIP-712 摘要 → ecrecover 还原签名地址 → 必须 == 请求 wallet（WALLET-04）
// → 消费 nonce 台账（一次性，WALLET-02 重放拒绝）。
func NewWalletAuth(nonceRepo *repository.WalletNonceRepository, chainID uint64, expectedAction string) gin.HandlerFunc {
	return func(c *gin.Context) {
		wallet := c.GetHeader(HeaderWalletAddr)
		nonce := c.GetHeader(HeaderWalletNonce)
		action := c.GetHeader(HeaderWalletAction)
		sigHex := c.GetHeader(HeaderWalletSig)

		if wallet == "" || nonce == "" || sigHex == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: missing wallet auth headers"})
			return
		}
		if !common.IsHexAddress(wallet) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: invalid wallet address"})
			return
		}
		// action 绑端点（WALLET-ACTION）：expectedAction 非空且必须与签名摘要里的 action 严格一致，
		// 否则为某动作签的 nonce 被挪用到别的端点（如 withdraw 的签名打 bills/pay）→ 拒绝。
		if expectedAction == "" || action != expectedAction {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: action mismatch"})
			return
		}

		recovered, err := recoverSigner(wallet, nonce, action, sigHex, chainID)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: " + err.Error()})
			return
		}
		// 绑定 msg.sender 语义：签名地址必须 == 请求 wallet（WALLET-04），否则拒绝。
		if !strings.EqualFold(recovered.Hex(), common.HexToAddress(wallet).Hex()) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: signature does not match wallet"})
			return
		}
		// 🔴 N1 一次性 nonce 台账消费（WALLET-02 重放拒绝）：未签发/已用/钱包不符/过期均失败。
		if !nonceRepo.Consume(wallet, nonce) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized: nonce invalid, used, or expired"})
			return
		}

		c.Set(CtxWallet, strings.ToLower(common.HexToAddress(wallet).Hex()))
		c.Next()
	}
}

// recoverSigner 重算 EIP-712 摘要并从签名还原签名地址。
func recoverSigner(wallet, nonce, action, sigHex string, chainID uint64) (common.Address, error) {
	digest, err := WalletAuthDigest(wallet, nonce, action, chainID)
	if err != nil {
		return common.Address{}, fmt.Errorf("digest: %w", err)
	}

	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		// 容错：允许不带 0x 前缀的纯十六进制。
		if s, derr := hexutil.Decode("0x" + strings.TrimPrefix(sigHex, "0x")); derr == nil {
			sig = s
		} else {
			return common.Address{}, fmt.Errorf("invalid signature encoding")
		}
	}
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length")
	}
	// 归一化 V：以太坊钱包常用 27/28，go-ethereum SigToPub 期望 0/1。
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pub, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("ecrecover failed")
	}
	return crypto.PubkeyToAddress(*pub), nil
}
