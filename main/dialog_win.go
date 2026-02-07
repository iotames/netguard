//go:build windows
// +build windows

package main

import (
	sqdialog "github.com/sqweek/dialog"
)

func errorMsg(title string, format string, args ...interface{}) {
	// sqdialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
	sqdialog.Message(format, args...).Title(title).Error()
}
