package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/iotames/netguard"
)

var Devname string
var ListDev bool

func main() {
	if len(os.Args) > 1 {
		arg1 := os.Args[1]
		versionArgs := []string{"--version", "-v", "-version"}
		if slices.Contains(versionArgs, arg1) {
			versionInfo()
			return
		}
	}

	// 设置程序日志
	f := setLog()
	defer f.Close()

	if ListDev {
		showDevices()
		return
	}

	// 设置geoip数据库
	setGeoipDb()

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
	netguard.Run(Devname)
}

func init() {
	flag.StringVar(&Devname, "devname", "", `netguard.exe --devname="\Device\NPF_{3757BF1E-96B9-441B-8D4B-95EAB49ECA36}"`)
	flag.BoolVar(&ListDev, "listdev", false, "netguard.exe --listdev")
	flag.Parse()
}
