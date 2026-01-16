package main

import (
	// "flag"
	"fmt"
	"os"

	// "log/slog"
	"github.com/iotames/netguard"
	"github.com/iotames/netguard/log"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "-version" {
			versionInfo()
			return
		}
	}

	// log.SetLevel(slog.LevelInfo)
	// f, err := log.SetLogWriterByFile("netguard.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	geoipFile := "GeoLite2-City.mmdb"
	err := netguard.SetGeoipDb(geoipFile)
	if err != nil {
		log.Error("set geoip db fail", "error", err.Error(), "geoipFile", geoipFile)
		panic(err)
	}

	var ipinfomap = make(map[string]string, 1000)
	netguard.SetPacketHook(func(info *netguard.TrafficRecord) {
		remoteIp := info.RemoteIP.String()
		remoteInfostr, ok := ipinfomap[remoteIp]
		if !ok {
			ipinfo := netguard.GetIpGeo(remoteIp)
			remoteInfostr = fmt.Sprintf("%s %s", ipinfo.Country, ipinfo.City)
			ipinfomap[remoteIp] = remoteInfostr
		}
		totalMB := float64(info.BytesReceived+info.BytesSent) / 1024.0 / 1024.0
		fmt.Printf("流量概要:%s, 流量:%.2fMB, IP解析:%s\n", info.Msg, totalMB, remoteInfostr)
	})
	netguard.Run()
}

// func init() {
// 	flag.Parse()
// }
