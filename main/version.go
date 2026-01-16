package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

var (
	BuildTime string
	Version   string
)

func versionInfo() {
	fmt.Printf("NetGuard:%s, BuildTime:%s\n", Version, BuildTime)
	// 方法1: 位运算法
	bit1 := 32 << (^uint(0) >> 63)
	fmt.Printf("当前系统位数 (位运算法): %d 位\n", bit1)

	// 方法2: unsafe.Sizeof法
	bit2 := unsafe.Sizeof(uint(0)) * 8
	fmt.Printf("当前系统位数 (Sizeof法): %d 位\n", bit2)

	// 附加信息：查看GOARCH环境变量（Go编译器目标架构）
	fmt.Printf("GOARCH 环境变量: %s\n", runtime.GOARCH)
}
