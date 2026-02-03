package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/iotames/netguard"
	"github.com/iotames/netguard/log"
)

func setLog() *os.File {
	// 设置日志文件
	log.SetLevel(slog.LevelInfo)
	logFile := "netguard.log"
	f, err := log.SetLogWriterByFile(logFile)
	if err != nil {
		log.Error("set log file fail", "error", err.Error(), "logFile", logFile)
		panic(err)
	}
	return f
	// defer f.Close()
}

func setGeoipDb() {
	geoipFile := "GeoLite2-City.mmdb"
	err := netguard.SetGeoipDb(geoipFile)
	if err != nil {
		log.Error("set geoip db fail", "error", err.Error(), "geoipFile", geoipFile)
		panic(err)
	}
}

func showDevices() {
	devs := netguard.GetDeviceList()
	for i, dev := range devs {
		fmt.Printf("---[%d]--Name(%s)--Description(%s)--------\n", i, dev.Name, dev.Description)
	}
}
