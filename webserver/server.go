package webserver

import (
	"fmt"
	"sync"

	e "github.com/iotames/easyserver"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/netguard"
	"github.com/iotames/netguard/db"
	"github.com/iotames/netguard/device"
	"github.com/iotames/netguard/hotswap"
)

var AppTitle = "NetGuard网络流量监控"

func Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	svr := e.NewServer(addr)
	setMiddlewares(svr)
	setHandler(svr)
	return svr.ListenAndServe()
}

func setMiddlewares(svr *httpsvr.EasyServer) {
	svr.AddMiddleHead(httpsvr.NewMiddleCORS("*"))
	svr.AddMiddleHead(httpsvr.NewMiddleStatic("/static", "./static"))
}

func setHandler(svr *httpsvr.EasyServer) {
	svr.AddHandler("GET", "/", home)
	svr.AddHandler("GET", "/debug", debug)
	svr.AddHandler("GET", "/api/log/setfile", setlogfile)
	svr.AddHandler("POST", "/api/uploadfile", uploadfile)
	svr.AddHandler("GET", "/api/device/list", deviceList)
	svr.AddHandler("GET", "/api/amis-page-config", getAmisPageConfig)
	svr.AddHandler("POST", "/api/netguard/start", netguardStart)
}

type NetguardConf struct {
	DevName string `json:"devname"`
}

var startConf NetguardConf
var netguardStarted bool

func netguardStart(ctx httpsvr.Context) {
	fmt.Printf("---startConf11(%+v)-----\n", startConf)
	if netguardStarted {
		e.ResponseJsonFail(ctx, "请先停止后启动", 500)
		return
	}

	err := ctx.GetPostJson(&startConf)
	if err != nil {
		e.ResponseJsonFail(ctx, err.Error(), 500)
		return
	}
	fmt.Printf("---startConf22(%+v)-----\n", startConf)
	go func() {
		netguardStarted = true

		d := db.GetDb()

		// 使用sync.Map替代map，避免出现concurrent map writes错误
		var ipinfomap = &sync.Map{}

		netguard.SetPacketHook(func(info *netguard.TrafficRecord) {
			remoteIp := info.RemoteIP.String()
			// 跳过本地IP的处理
			if netguard.IsNativeIP(remoteIp) {
				return
			}
			var ipinfo netguard.GeoIpInfo
			// 使用sync.Map的方法
			if remoteInfo, ok := ipinfomap.Load(remoteIp); ok {
				ipinfo = remoteInfo.(netguard.GeoIpInfo)
				// totalMB := float64(info.BytesReceived+info.BytesSent) / 1024.0 / 1024.0
				// logmsg := fmt.Sprintf("流量概要:%s, 流量:%.2fMB, IP解析:%s", info.Msg, totalMB, remoteInfostr)
				// fmt.Printf("%s\n", logmsg)
				// log.Info("PacketHook", "logmsg", logmsg)

			} else {
				ipinfo = netguard.GetIpGeo(remoteIp)
				// remoteInfostr := fmt.Sprintf("%s %s", ipinfo.Country, ipinfo.City)
				ipinfomap.Store(remoteIp, ipinfo)
				// totalMB := float64(info.BytesReceived+info.BytesSent) / 1024.0 / 1024.0
				// logmsg := fmt.Sprintf("流量概要:%s, 流量:%.2fMB, IP解析:%s", info.Msg, totalMB, remoteInfostr)
				// fmt.Printf("%s\n", logmsg)
				// log.Info("PacketHook", "logmsg", logmsg)

			}
			_, err = d.Exec(`INSERT INTO ng_hook_logs (
            remote_ip, remote_port, protocol,
			process_name, process_pid,
            bytes_current_len, inbound,
            ip_country, ip_city
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				remoteIp, info.RemotePort, info.Protocol,
				info.ProcessName, info.ProcessPID,
				info.BytesCurrentLen, info.Inbound,
				ipinfo.Country, ipinfo.City,
			)

			if err != nil {
				fmt.Println("sql error:", err.Error())
			}

		})

		netguard.Run(startConf.DevName)

	}()
	e.ResponseJsonOk(ctx, "启动成功")
}

func setlogfile(ctx httpsvr.Context) {
	err := setLogFile()
	if err != nil {
		e.ResponseJsonFail(ctx, err.Error(), 500)
		return
	}
	e.ResponseJsonOk(ctx, "设置成功")
}

func deviceList(ctx httpsvr.Context) {
	devlist := device.GetDeviceList()

	options := make([]map[string]string, len(devlist))
	for i, v := range devlist {
		options[i] = map[string]string{"label": v.Description, "value": v.Name}
	}
	// json返回
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"options": options}, "success", 0).Bytes())
}

func home(ctx httpsvr.Context) {
	data := map[string]interface{}{
		"title": AppTitle,
		// "web_server_port": webServerPort,
	}
	SetContentByTplFile("amis.html", ctx.Writer, data)
	// ctx.Writer.Write(response.NewApiDataOk("hello api").Bytes())
}

func debug(ctx httpsvr.Context) {
	sd := hotswap.GetScriptDir(nil)
	stxt, err := sd.GetScriptText("init.sql")
	if err != nil {
		e.ResponseJsonFail(ctx, err.Error(), 500)
		return
	}
	ctx.Writer.Write(response.NewApiDataOk(stxt).Bytes())
}
