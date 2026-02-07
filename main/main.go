package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/iotames/netguard"
	"github.com/iotames/netguard/conf"
	"github.com/iotames/netguard/webserver"
)

func main() {
	var err error
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
		// f := setLog()
		// defer f.Close()
		err = webserver.Run(Port)
		if err != nil {
			panic(fmt.Errorf("webserver.Run err(%v)", err))
		}
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
			errorMsg("初始化错误", "配置初始化错误。请检查目录权限，或尝试以管理员身份运行程序。conf.LoadEnv err:%s", err.Error())
		}
		panic(fmt.Errorf("init err(%v)", err))
	}
	parseArgs()
	initScript()
}
