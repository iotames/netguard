package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iotames/netguard"
)

var Devname string
var ListDev bool

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "-version" {
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

	var ipinfomap = make(map[string]string, 1000)
	fmt.Println("SetPacketHook Start")
	netguard.SetPacketHook(func(info *netguard.TrafficRecord) {
		fmt.Println("PacketHook Start")
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
	netguard.Run(Devname)
}

func init() {
	flag.StringVar(&Devname, "devname", "", `netguard.exe --devname="\Device\NPF_{3757BF1E-96B9-441B-8D4B-95EAB49ECA36}"`)
	flag.BoolVar(&ListDev, "listdev", false, "netguard.exe --listdev")
	flag.Parse()
}
