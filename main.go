package main

import (
	"fmt"
	"interpreter/repl"
	"os"
)

func main() {
	fmt.Println("哈语言 — 只有哈和空格")
	fmt.Println("1-8个哈 = 操作，空格 = 参数")
	fmt.Println("空行 = 执行程序")
	fmt.Println()
	repl.Start(os.Stdin, os.Stdout)
}
