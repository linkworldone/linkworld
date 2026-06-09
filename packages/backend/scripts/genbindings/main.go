// genbindings 用 go-ethereum 的 accounts/abi/bind 库直接生成 abigen 风格的合约绑定。
//
// 为什么不用 `abigen` 命令：go-ethereum v1.13.5 的 cmd/abigen 依赖
// github.com/fjl/memsize，该包在 Go 1.21+ 下因 runtime.stopTheWorld 符号变更
// 无法链接（link: invalid reference to runtime.stopTheWorld）。bind.Bind 库本身
// 不受影响，且生成结果与 abigen 命令完全一致（abigen 内部就是调用 bind.Bind）。
//
// 用法（一般通过 scripts/gen-bindings.sh 调用）：
//
//	go run ./scripts/genbindings -abi <file.json> -pkg bindings -type <Name> -out <out.go>
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func main() {
	abiPath := flag.String("abi", "", "path to JSON file containing the contract ABI array")
	pkg := flag.String("pkg", "bindings", "Go package name for the generated file")
	typeName := flag.String("type", "", "Go type name for the binding (contract name)")
	out := flag.String("out", "", "output .go file path")
	flag.Parse()

	if *abiPath == "" || *typeName == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: genbindings -abi <file> -type <Name> -out <file> [-pkg bindings]")
		os.Exit(2)
	}

	abiBytes, err := os.ReadFile(*abiPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read abi %s: %v\n", *abiPath, err)
		os.Exit(1)
	}
	abiStr := strings.TrimSpace(string(abiBytes))

	// bind.Bind 接受多组并行切片；本工具一次只绑一个合约。bytecode/sigs/aliases
	// 为空即可（仅生成只读/调用绑定，不内联部署字节码）。
	code, err := bind.Bind(
		[]string{*typeName},
		[]string{abiStr},
		[]string{""},        // bytecodes
		nil,                 // fsigs
		*pkg,                // package name
		bind.LangGo,         // target language
		nil,                 // libs
		nil,                 // aliases
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bind %s: %v\n", *typeName, err)
		os.Exit(1)
	}

	if err := os.WriteFile(*out, []byte(code), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", *out, err)
		os.Exit(1)
	}
	fmt.Printf("generated %s -> %s\n", *typeName, *out)
}
