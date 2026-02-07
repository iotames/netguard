package device

import (
	"fmt"
	"strings"

	"github.com/google/gopacket/pcap"
	// "github.com/iotames/netguard/log"
)

func GetDeviceList() []pcap.Interface {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		// log.Error("无法获取网络设备列表:", "错误", err)
		panic(fmt.Errorf("无法获取网络设备列表: %w", err))
	}
	return devices
}

// GetDefaultDevice 获取默认的网络设备。选择第一个非环回接口。
func GetDefaultDevice() pcap.Interface {
	// fmt.Printf("------------GetDefaultDevice: pcap.FindAllDevs() Start-----\n")
	// log.Info("GetDefaultDevice: pcap.FindAllDevs() Start")
	devices := GetDeviceList()

	// fmt.Printf("----------GetDefaultDevice: pcap.FindAllDevs() Done--设备数(%d)---\n", len(devices))
	// log.Info("GetDefaultDevice: pcap.FindAllDevs() Done", "设备数", len(devices))

	var defaultDevice pcap.Interface
	for _, device := range devices {
		// fmt.Printf("------GetDefaultDevice device Info:---设备名(%+v)--设备描述(%+v)---设备地址(%+v)---\n", device.Name, device.Description, device.Addresses)
		// log.Info("GetDefaultDevice device Info:", "设备名", device.Name, "设备描述", device.Description, "设备地址", device.Addresses)
		// 跳过虚拟和蓝牙设备
		descLower := strings.ToLower(device.Description)
		if strings.Contains(descLower, "bluetooth") ||
			strings.Contains(descLower, "virtual") ||
			strings.Contains(descLower, "loopback") {
			continue
		}
		// 选择第一个有IP地址的非环回设备
		if len(device.Addresses) > 0 {
			for _, addr := range device.Addresses {
				if addr.IP != nil && !addr.IP.IsLoopback() {
					defaultDevice = device
					break
				}
			}
		}
		if defaultDevice.Name != "" {
			break
		}
	}
	if defaultDevice.Name == "" {
		// log.Error("未找到合适的网络设备")
		panic("未找到合适的网络设备")
	}
	return defaultDevice
}
