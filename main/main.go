package main

import (
	// "flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"

	"github.com/iotames/netguard"
	"github.com/iotames/netguard/log"
)

func main() {
	if len(os.Args) > 1 {
		arg1 := os.Args[1]
		versionArgs := []string{"--version", "-v", "-version"}
		if slices.Contains(versionArgs, arg1) {
			versionInfo()
			return
		}
		if arg1 == "log" {
			logstart()
		}
	}

	geoipFile := "GeoLite2-City.mmdb"
	err := netguard.SetGeoipDb(geoipFile)
	if err != nil {
		log.Error("set geoip db fail", "error", err.Error(), "geoipFile", geoipFile)
		panic(err)
	}
	// 使用sync.Map替代map，避免出现concurrent map writes错误
	var ipinfomap = &sync.Map{}

	netguard.SetPacketHook(func(info *netguard.TrafficRecord) {
		remoteIp := info.RemoteIP.String()
		// 跳过本地IP的处理
		if netguard.IsNativeIP(remoteIp) {
			return
		}
		// 使用sync.Map的方法
		if remoteInfostr, ok := ipinfomap.Load(remoteIp); ok {
			totalMB := float64(info.BytesReceived+info.BytesSent) / 1024.0 / 1024.0
			fmt.Printf("流量概要:%s, 流量:%.2fMB, IP解析:%s\n", info.Msg, totalMB, remoteInfostr)
		} else {
			ipinfo := netguard.GetIpGeo(remoteIp)
			remoteInfostr := fmt.Sprintf("%s %s", ipinfo.Country, ipinfo.City)
			ipinfomap.Store(remoteIp, remoteInfostr)

			totalMB := float64(info.BytesReceived+info.BytesSent) / 1024.0 / 1024.0
			fmt.Printf("流量概要:%s, 流量:%.2fMB, IP解析:%s\n", info.Msg, totalMB, remoteInfostr)
		}
	})
	netguard.Run()
}

func logstart() {
	log.SetLevel(slog.LevelInfo)
	f, err := log.SetLogWriterByFile("netguard.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

// func init() {
// 	flag.Parse()
// }
