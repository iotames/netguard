package conf

import (
	"fmt"
	"os"

	"github.com/iotames/easyconf"
)

var cf *easyconf.Conf

const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"

const DEFAULT_ENV_FILE = ".env"
const DEFAULT_RUNTIME_DIR = "runtime"
const DEFAULT_RESOURCE_DIR = "resource"
const DEFAULT_SCRIPTS_DIR = "scripts"
const DEFAULT_WEB_SERVER_PORT = 8080
const DEFAULT_DB_DRIVER = DRIVER_POSTGRES
const DEFAULT_DB_HOST = "127.0.0.1"
const DEFAULT_DB_PORT = 5432
const DEFAULT_DB_NAME = "postgres"
const DEFAULT_DB_SCHEMA = "public"
const DEFAULT_DB_USERNAME = "postgres"
const DEFAULT_DB_PASSWORD = "postgres"

var RuntimeDir string

// var ResourceDir string
var ScriptsDir string

var WebServerPort int
var ShowSql bool
var DbDriver, DbHost, DbName, DbSchema, DbUsername, DbPassword string
var DbPort int

func getEnvFile() string {
	efile := os.Getenv("NGD_ENV_FILE")
	if efile == "" {
		efile = DEFAULT_ENV_FILE
	}
	return efile
}

func setConfByEnv() {
	// # 设置 NGD_ENV_FILE 环境变量，可更改配置文件路径。
	cf = easyconf.NewConf(getEnvFile())

	// cf.StringVar(&ResourceDir, "RESOURCE_DIR", DEFAULT_RESOURCE_DIR, "")
	cf.StringVar(&RuntimeDir, "RUNTIME_DIR", DEFAULT_RUNTIME_DIR, "")
	cf.StringVar(&ScriptsDir, "SCRIPTS_DIR", DEFAULT_SCRIPTS_DIR, "放自定义的脚本文件")
	cf.IntVar(&WebServerPort, "WEB_SERVER_PORT", DEFAULT_WEB_SERVER_PORT, "启动Web服务器的端口号")

	cf.BoolVar(&ShowSql, "SHOW_SQL", false, "是否输出SQL调试信息")
	cf.StringVar(&DbDriver, "DB_DRIVER", DEFAULT_DB_DRIVER, "数据库类型: mysql,sqlite3,postgres")
	cf.StringVar(&DbHost, "DB_HOST", DEFAULT_DB_HOST, "数据库主机地址")
	cf.StringVar(&DbName, "DB_NAME", DEFAULT_DB_NAME, "数据库名")
	cf.StringVar(&DbSchema, "DB_SCHEMA", DEFAULT_DB_SCHEMA, "数据库schema")
	cf.IntVar(&DbPort, "DB_PORT", DEFAULT_DB_PORT, "数据库端口号:5432(postgres);3306(mysql)")
	cf.StringVar(&DbUsername, "DB_USERNAME", DEFAULT_DB_USERNAME, "数据库用户名")
	cf.StringVar(&DbPassword, "DB_PASSWORD", DEFAULT_DB_PASSWORD, "数据库密码")

	cf.Parse(false)
}

func ShowConf() {
	fmt.Println(cf.String())
}

func LoadEnv() error {
	var err error
	setConfByEnv()
	// if !IsPathExists(ResourceDir) {
	// 	err = Mkdir(ResourceDir)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	if !IsPathExists(RuntimeDir) {
		fmt.Printf("------创建runtime目录(%s)--\n", RuntimeDir)
		err = os.Mkdir(RuntimeDir, 0755)
		if err != nil {
			fmt.Printf("----runtime目录(%s)创建失败(%v)---\n", RuntimeDir, err)
		}
	}
	return err
}

// IsPathExists 判断文件或文件夹是否存在
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// fmt.Println(stat.IsDir())
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Mkdir 创建目录
func Mkdir(path string) error {
	if IsPathExists(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
