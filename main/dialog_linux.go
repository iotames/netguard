//go:build linux
// +build linux

package main

import (
	"fmt"
)

func errorMsg(title string, format string, args ...interface{}) {
	fmt.Println(title)
}
