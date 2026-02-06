package webserver

import (
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/netguard/device"
	"github.com/iotames/netguard/webserver/amis"
)

func getAmisPageConfig(ctx httpsvr.Context) {
	defaultDev := device.GetDefaultDevice()
	pageConf := amis.NewPage(AppTitle)
	item1 := amis.NewFormItem().Set("label", "监控网卡").Set("type", "select").Set("name", "devname").Set("value", defaultDev.Name).Set("source", "/api/device/list")
	// item2 := amis.NewFormItem().Set("type", "input-file").Set("name", "inputfile").Set("accept", ".xlsx").Set("label", "上传.xlsx文件").Set("maxSize", 10048576).Set("receiver", "/api/uploadfile")
	pageConf.Body = *amis.NewForm("/api/netguard/start").AddItem(item1).SetSubmitText("启动")
	// .SetTitle("AppTitle")
	// .AddItem(item2)
	ctx.Writer.Write(response.NewApiData(pageConf.Json(), "success", 0).Bytes())
}
