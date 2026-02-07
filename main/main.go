package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/iotames/netguard"
	"github.com/iotames/netguard/conf"
	"github.com/iotames/netguard/webserver"
	sqdialog "github.com/sqweek/dialog"
)

func main() {
	if V || VersionV {
		versionInfo()
		return
	}
	if ListDev {
		showDevices()
		return
	}

	if runtime.GOOS == "windows" && IsPathExists("amis.html") && Port > 0 {
		go func() {
			time.Sleep(1 * time.Second)
			err := StartBrowserByUrl(fmt.Sprintf("http://127.0.0.1:%d", Port))
			if err != nil {
				println(err)
			}
		}()
	}
	if Port > 0 {
		webserver.Run(Port)
	} else {
		runNetguard()
	}
}

func runNetguard() {
	// 设置程序日志
	f := setLog()
	defer f.Close()

	// 设置geoip数据库
	setGeoipDb()
	netguard.DebugRun(Devname)
}

func init() {
	err := conf.LoadEnv()
	if err != nil {
		if runtime.GOOS == "windows" {
			// sqdialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
			sqdialog.Message("环境变量初始化错误（conf.LoadEnv err）:%s", err.Error()).Title("初始化错误").Error()
		}
		panic(fmt.Errorf("init err(%v)", err))
	}
	parseArgs()
	initScript()
}
