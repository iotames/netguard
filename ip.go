package netguard

import (
	"fmt"
	"net"
	"sync"

	"github.com/iotames/netguard/log"
	gnet "github.com/shirou/gopsutil/v3/net"
)

// 全局配置
var (
	realTimeProcessQuery = true                   // 实时进程查询开关
	processQueryCache    = make(map[string]int32) // 进程查询缓存
	processCacheMutex    sync.RWMutex
)

// updateLocalIPs 获取本机所有IP地址
func updateLocalIPs() {
	var ips []net.IP

	// 获取所有网络接口
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error("获取网络接口失败:", "错误", err)
		return
	}

	for _, iface := range ifaces {
		// 跳过环回和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() {
				ips = append(ips, ip)
			}
		}
	}

	// 使用写锁替换整个切片，避免并发读取时竞争
	localIPsMutex.Lock()
	localIPs = ips
	localIPsMutex.Unlock()
}

// findPidByConnection 通过IP和端口查找对应的进程PID
//
// 数据包到达
//
//	↓
//
// processCapturedPacket调用findPIDByConnection
//
//	↓
//
// findPIDByConnection：只查缓存和映射表，返回PID（可能为0）
//
//	↓
//
// 调用updateTrafficStats
//
//	↓
//
// 如果是新建连接 && 开启实时查询 && PID=0 → 触发实时查询
//
//	↓
//
// 更新流量统计
func findPidByConnection(ip net.IP, port uint16) int32 {
	key := fmt.Sprintf("%s:%d", ip.String(), port)

	// 1. 先查缓存（快速路径）
	processCacheMutex.RLock()
	if pid, exists := processQueryCache[key]; exists {
		processCacheMutex.RUnlock()
		return pid
	}
	processCacheMutex.RUnlock()

	// 2. 查全局连接映射表
	if pid, exists := connectionMap.Load(key); exists {
		if pidInt, ok := pid.(int32); ok && pidInt > 0 {
			// 更新缓存
			processCacheMutex.Lock()
			processQueryCache[key] = pidInt
			processCacheMutex.Unlock()
			return pidInt
		}
	}

	// // 3. 如果开启实时查询且缓存未命中，进行实时查询
	// if realTimeProcessQuery {
	// 	return queryProcessRealTime(ip, port)
	// }

	return 0
}

// queryProcessRealTime 实时查询进程信息
func queryProcessRealTime(ip net.IP, port uint16) int32 {
	// 立即查询当前系统连接表
	connections, err := gnet.Connections("all")
	if err != nil {
		log.Warn("获取网络连接信息失败:", "错误", err)
		return 0
	}

	targetIP := ip.String()
	for _, conn := range connections {
		if conn.Laddr.IP == targetIP && conn.Laddr.Port == uint32(port) {
			// 找到匹配连接，更新缓存和全局映射
			key := fmt.Sprintf("%s:%d", ip.String(), port)

			processCacheMutex.Lock()
			processQueryCache[key] = conn.Pid
			processCacheMutex.Unlock()

			connectionMap.Store(key, conn.Pid)
			log.Debug("实时查询PID成功", "key", key, "PID", conn.Pid)
			return conn.Pid
		}
	}

	return 0
}

// isLocalIP 判断一个IP地址是否为本地IP
func isLocalIP(ip net.IP) bool {
	// 使用读锁保护 localIPs 访问
	localIPsMutex.RLock()
	defer localIPsMutex.RUnlock()
	for _, localIP := range localIPs {
		if localIP.Equal(ip) {
			return true
		}
	}
	return false
}

// IsNativeIP 判断给定的IP字符串是否属于"本地"地址范畴。
// 本地地址是指那些通常不需要进行外部地理定位查询的IP地址，包括：
// 1. 组播地址（Multicast）：用于一对多通信，无实际地理位置
// 2. 私有地址（Private/Internal）：在局域网内使用，不在公网路由
// 3. 特殊保留地址（Reserved/Special）：用于协议或系统功能
// 这类地址的共同特点是：它们不代表互联网上的真实公网主机位置。
//
// 参数：
//
//	ipStr - 要检查的IP地址字符串，如"192.168.1.1"或"239.255.255.250"
//
// 返回值：
//
//	bool - 如果是本地地址返回true，如果是公网地址返回false
//
// 设计说明：
//
//	函数采用分层检查策略，按照地址类别从最常见到较少见顺序检查，
//	以提高平均执行效率。当前版本主要处理IPv4地址。
func IsNativeIP(ipStr string) bool {
	// 第一步：将字符串解析为net.IP对象
	// net.ParseIP能识别IPv4和IPv6格式，并返回规范化后的字节表示
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// 无法解析的字符串格式，不是有效的IP地址
		// 为了系统健壮性，将其视为本地地址跳过处理
		return true
	}

	// 第二步：检查是否为IPv4地址
	// To4()方法将IPv4地址返回4字节表示，IPv6地址返回nil
	// 这种方式能正确处理IPv4映射的IPv6地址(::ffff:192.168.1.1)
	if ipv4 := ip.To4(); ipv4 != nil {
		// 获取IPv4地址的第一个字节（8位）
		// IPv4地址如A.B.C.D，firstOctet对应A的值
		firstOctet := ipv4[0]

		// ---------- 1. 检查组播地址 (224.0.0.0 - 239.255.255.255) ----------
		// 组播地址是D类地址，第一个字节在224-239范围内
		// 这些地址用于一对多通信，没有实际地理位置
		if firstOctet >= 224 && firstOctet <= 239 {
			return true // 是组播地址，属于本地地址
		}

		// ---------- 2. 检查各类私有和保留地址 ----------
		// 使用switch-case结构清晰表达多个独立条件

		switch {
		// RFC 1918定义的A类私有地址范围：10.0.0.0/8
		// 第一个字节为10，后面三个字节任意值都是私有地址
		case firstOctet == 10:
			return true // 10.x.x.x

		// RFC 1918定义的B类私有地址范围：172.16.0.0/12
		// 第一个字节为172，第二个字节在16-31之间
		// 注意：172.16.x.x - 172.31.x.x是私有地址，但172.15.x.x和172.32.x.x不是
		case firstOctet == 172 && ipv4[1] >= 16 && ipv4[1] <= 31:
			return true // 172.16.x.x - 172.31.x.x

		// RFC 1918定义的C类私有地址范围：192.168.0.0/16
		// 第一个字节为192，第二个字节为168
		case firstOctet == 192 && ipv4[1] == 168:
			return true // 192.168.x.x

		// RFC 1122定义的环回地址范围：127.0.0.0/8
		// 用于本机内部通信，所有127.x.x.x地址都指向本机
		case firstOctet == 127:
			return true // 127.x.x.x

		// RFC 3927定义的链路本地地址范围：169.254.0.0/16
		// 当DHCP失败时系统自动分配，仅在本地链路有效
		case firstOctet == 169 && ipv4[1] == 254:
			return true // 169.254.x.x

		// RFC 1122定义的"本网络"地址：0.0.0.0/8
		// 通常用作默认路由或表示无效地址
		case firstOctet == 0:
			return true // 0.x.x.x
		}

		// ---------- 3. 其他可能的本地地址检查 ----------
		// 以下是不太常见但仍可能遇到的本地地址，按需启用：

		// RFC 6598定义的共享地址空间：100.64.0.0/10
		// 用于运营商级NAT场景 (100.64.0.0 - 100.127.255.255)
		// if firstOctet == 100 && (ipv4[1] & 0xC0) == 64 {
		//     return true
		// }

		// RFC 6890定义的文档/测试地址：192.0.2.0/24, 198.51.100.0/24, 203.0.113.0/24
		// 用于文档和示例代码，不应出现在公网
		// if (firstOctet == 192 && ipv4[1] == 0 && ipv4[2] == 2) ||
		//    (firstOctet == 198 && ipv4[1] == 51 && ipv4[2] == 100) ||
		//    (firstOctet == 203 && ipv4[1] == 0 && ipv4[2] == 113) {
		//     return true
		// }

		// 通过所有检查，这是一个公网IPv4地址
		return false
	}

	// 第三步：可选的IPv6地址检查
	// 当前版本暂不实现IPv6检查，可根据需要扩展
	// IPv6的本地地址包括：
	//   - 组播地址: ff00::/8
	//   - 链路本地地址: fe80::/10
	//   - 唯一本地地址: fc00::/7 (类似IPv4的私有地址)
	//   - 环回地址: ::1

	// 默认情况下，暂时不将IPv6地址视为本地地址
	// 因为IPv6地理定位可能仍有价值，且IPv6地址空间巨大
	return false
}
