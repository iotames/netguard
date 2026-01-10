package main

import (
	"fmt"

	// "log/slog"
	// "github.com/iotames/netguard/log"
	"github.com/iotames/netguard"
)

func main() {
	// log.SetLevel(slog.LevelInfo)
	// f, err := log.SetLogWriterByFile("netguard.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	netguard.SetPacketHook(func(info *netguard.TrafficRecord) {
		fmt.Println("PacketInfo:", info.Msg, "Local:", info.LocalIP, info.LocalPort)

		// 字节转MB，保留两位小数
		// currentMB := float64(tr.BytesCurrentLen) / 1024.0 / 1024.0
		// totalMB := float64(tr.BytesReceived+tr.BytesSent) / 1024.0 / 1024.0
		// tr.Msg = fmt.Sprintf("%s-%s, Remote(%s:%d), Process(%d-%s), Length(%.2fMB/%.2fMB)", tr.Protocol, direction, remoteIP.String(), remotePort, tr.ProcessPID, tr.ProcessName, currentMB, totalMB)
	})

	// for now := range time.Tick(10 * time.Second) {
	// 	fmt.Println("当前时间：", now, "IpMapLength:", len(netguard.GetTrafficStats()))
	// }

	netguard.Run()
}
