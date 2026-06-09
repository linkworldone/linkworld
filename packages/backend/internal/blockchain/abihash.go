// Package blockchain 的 abihash 工具：对合约 ABI 计算指纹，与 deployments/<net>.json
// 的 abiHash 字段比对，确认后端绑定与链上部署的合约 ABI 同步。
//
// 算法必须与合约侧 deploy.ts 完全一致：
//
//	abiHash = keccak256( utf8( ethers.Interface(abi).formatJson() ) )
//
// 关键点：ethers v6 的 `Interface.formatJson()` 输出的不是 hardhat artifact 原始
// .abi 数组，而是一份**规范化**的紧凑 JSON——字段顺序固定、补齐 constant/payable
// 布尔、nonpayable 时省略 stateMutability。直接对原始 artifact .abi 取 keccak 不匹配
// （已实测：8 个合约原始 JSON 哈希全部对不上，formatJson 全部对上）。
//
// 本文件用纯 Go 复刻 formatJson 的规范化规则（见 canonicalizeABI），不引入 JS 运行时，
// 保证 go test 离线可跑。
package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// abiParam 对应 ethers formatJson 里 inputs/outputs/components 的单个参数项。
// 字段顺序即 ethers 的输出顺序：type, name, [indexed], [components]。
// 用指针 / omitempty 控制可选字段的存在性。
type abiParam struct {
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	Indexed    *bool      `json:"indexed,omitempty"`
	Components []abiParam `json:"components,omitempty"`
}

// rawParam 对应 hardhat artifact 里参数的原始字段（含 internalType，formatJson 会丢弃）。
type rawParam struct {
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	Indexed    *bool      `json:"indexed,omitempty"`
	Components []rawParam `json:"components,omitempty"`
}

// rawEntry 对应 hardhat artifact .abi 数组里的单个条目。
type rawEntry struct {
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	Anonymous       *bool      `json:"anonymous,omitempty"`
	StateMutability string     `json:"stateMutability,omitempty"`
	Inputs          []rawParam `json:"inputs,omitempty"`
	Outputs         []rawParam `json:"outputs,omitempty"`
}

// canonParam 在 marshal 时复刻 ethers 的参数字段顺序。
func canonParam(p rawParam, keepIndexed bool) abiParam {
	out := abiParam{Type: p.Type, Name: p.Name}
	if keepIndexed && p.Indexed != nil {
		out.Indexed = p.Indexed
	}
	if len(p.Components) > 0 {
		out.Components = make([]abiParam, len(p.Components))
		for i, c := range p.Components {
			// tuple 内部组件不带 indexed。
			out.Components[i] = canonParam(c, false)
		}
	}
	return out
}

func canonParams(ps []rawParam, keepIndexed bool) []abiParam {
	// ethers formatJson 总是输出 inputs/outputs 数组（即使为空 []），故不能用 nil。
	out := make([]abiParam, 0, len(ps))
	for _, p := range ps {
		out = append(out, canonParam(p, keepIndexed))
	}
	return out
}

// canonEntry 是 marshal 中转结构。各条目类型字段顺序固定，靠 RawMessage 拼装实现
// 「nonpayable 省略 stateMutability、view/pure/payable 保留」这类 ethers 特有规则。
//
// 为精确控制字段顺序与可选性，这里直接手工拼 JSON 片段而非依赖 struct tag——
// 因为 function 条目的字段顺序是 type,name,constant,[stateMutability],payable,inputs,outputs，
// 其中 stateMutability 的存在性取决于 mutability，struct tag 无法表达。
func canonicalizeEntryJSON(e rawEntry) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	writeField := func(b *bytes.Buffer, first bool, key string, val any) error {
		if !first {
			b.WriteByte(',')
		}
		kb, err := json.Marshal(key)
		if err != nil {
			return err
		}
		b.Write(kb)
		b.WriteByte(':')
		var vb bytes.Buffer
		ve := json.NewEncoder(&vb)
		ve.SetEscapeHTML(false)
		if err := ve.Encode(val); err != nil {
			return err
		}
		// Encode 追加了换行，去掉。
		b.Write(bytes.TrimRight(vb.Bytes(), "\n"))
		return nil
	}

	switch e.Type {
	case "function":
		buf.WriteByte('{')
		_ = writeField(&buf, true, "type", "function")
		_ = writeField(&buf, false, "name", e.Name)
		isView := e.StateMutability == "view" || e.StateMutability == "pure"
		isPayable := e.StateMutability == "payable"
		_ = writeField(&buf, false, "constant", isView)
		// stateMutability 仅在 view/pure/payable 时出现；nonpayable 省略。
		if e.StateMutability != "" && e.StateMutability != "nonpayable" {
			_ = writeField(&buf, false, "stateMutability", e.StateMutability)
		}
		_ = writeField(&buf, false, "payable", isPayable)
		_ = writeField(&buf, false, "inputs", canonParams(e.Inputs, false))
		_ = writeField(&buf, false, "outputs", canonParams(e.Outputs, false))
		buf.WriteByte('}')

	case "constructor":
		buf.WriteByte('{')
		_ = writeField(&buf, true, "type", "constructor")
		isPayable := e.StateMutability == "payable"
		_ = writeField(&buf, false, "constant", false)
		if e.StateMutability == "payable" {
			_ = writeField(&buf, false, "stateMutability", e.StateMutability)
		}
		_ = writeField(&buf, false, "payable", isPayable)
		_ = writeField(&buf, false, "inputs", canonParams(e.Inputs, false))
		buf.WriteByte('}')

	case "fallback", "receive":
		buf.WriteByte('{')
		_ = writeField(&buf, true, "type", e.Type)
		isPayable := e.StateMutability == "payable"
		_ = writeField(&buf, false, "constant", false)
		if e.StateMutability == "payable" {
			_ = writeField(&buf, false, "stateMutability", e.StateMutability)
		}
		_ = writeField(&buf, false, "payable", isPayable)
		buf.WriteByte('}')

	case "event":
		anon := false
		if e.Anonymous != nil {
			anon = *e.Anonymous
		}
		buf.WriteByte('{')
		_ = writeField(&buf, true, "type", "event")
		_ = writeField(&buf, false, "anonymous", anon)
		_ = writeField(&buf, false, "name", e.Name)
		_ = writeField(&buf, false, "inputs", canonParams(e.Inputs, true))
		buf.WriteByte('}')

	case "error":
		buf.WriteByte('{')
		_ = writeField(&buf, true, "type", "error")
		_ = writeField(&buf, false, "name", e.Name)
		_ = writeField(&buf, false, "inputs", canonParams(e.Inputs, false))
		buf.WriteByte('}')

	default:
		return nil, fmt.Errorf("abihash: unsupported ABI entry type %q", e.Type)
	}

	_ = enc // enc 仅占位，保留以备扩展
	return buf.Bytes(), nil
}

// CanonicalABIJSON 把 hardhat artifact 的原始 .abi 数组 JSON 转换成与 ethers v6
// Interface.formatJson() 字节级一致的规范 JSON。
func CanonicalABIJSON(rawABI []byte) ([]byte, error) {
	var entries []rawEntry
	if err := json.Unmarshal(rawABI, &entries); err != nil {
		return nil, fmt.Errorf("abihash: parse raw abi: %w", err)
	}
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, e := range entries {
		if i > 0 {
			buf.WriteByte(',')
		}
		b, err := canonicalizeEntryJSON(e)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// ComputeABIHash 对原始 ABI 数组 JSON 计算 abiHash（= keccak256(utf8(canonical formatJson))），
// 返回 0x 前缀十六进制字符串，与 deployments/<net>.json 的 abiHash 值同格式。
func ComputeABIHash(rawABI []byte) (string, error) {
	canon, err := CanonicalABIJSON(rawABI)
	if err != nil {
		return "", err
	}
	h := crypto.Keccak256(canon)
	return hexutil.Encode(h), nil
}

// VerifyABIHash 比对单个合约的实际 abiHash 与期望值（来自 deployments json）。
// 不一致返回带差异信息的 error。
func VerifyABIHash(name string, rawABI []byte, expected string) error {
	got, err := ComputeABIHash(rawABI)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if got != expected {
		return fmt.Errorf("abihash mismatch for %s: got %s, expected %s (ABI 与部署不同步，需重新生成绑定或核对合约)", name, got, expected)
	}
	return nil
}
