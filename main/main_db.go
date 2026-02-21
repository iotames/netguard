package main

import (
	"database/sql"
	"sync"
	"time"

	"github.com/iotames/easydb"
	"github.com/iotames/netguard/conf"
	"github.com/iotames/netguard/db"
	"github.com/iotames/netguard/hotswap"
	"github.com/iotames/netguard/log"

	// _ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var once sync.Once
var edb *easydb.EasyDb

func dbinit() {
	edb = newDb(conf.DbDriver, conf.DbHost, conf.DbUsername, conf.DbPassword, conf.DbName, conf.DbPort)
	_, err := createTables()
	if err != nil {
		panic(err)
	}
	db.SetDb(edb)
	log.Info("数据库初始化完成", "DbDriver", conf.DbDriver, "DbHost", conf.DbHost, "DbPort", conf.DbPort, "DbName", conf.DbName)

}

func createTables() (sql.Result, error) {
	if conf.DbDriver == conf.DRIVER_SQLITE {
		return execSqlFile("sqlite_init.sql")
	}
	return execSqlFile("sqlite_init.sql")
}

func execSqlFile(filename string, args ...any) (sql.Result, error) {
	var err error
	var sqltxt string
	sqltxt, err = hotswap.GetScriptDir(nil).GetSQL(filename)
	if err != nil {
		return nil, err
	}
	return edb.Exec(sqltxt, args...)
}

// DB结构体和方法，只给main,model调用
func newDb(driverName, dbHost, dbUser, dbPassword, dbName string, dbPort int) *easydb.EasyDb {
	var err error
	cf := easydb.NewDsnConf(driverName, dbHost, dbUser, dbPassword, dbName, dbPort)
	edb = easydb.NewEasyDbByConf(*cf)
	// 测试连接d
	if err = edb.Ping(); err != nil {
		panic(err)
	}
	// 设置合理的连接池参数
	// SHOW max_connections; // 检查 PostgreSQL 最大连接数
	// SELECT * FROM pg_stat_activity; 检查是否有其他应用占用连接
	// SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state = 'idle' AND query_start < NOW() - INTERVAL '10 minutes'; // 杀死空闲或长时间运行的连接
	// db.SetMaxOpenConns(20)  // 最大打开连接数
	// db.SetMaxIdleConns(5)   // 最大空闲连接数
	// db.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
	dbb := edb.GetSqlDB()
	dbb.SetMaxOpenConns(2000)
	dbb.SetConnMaxLifetime(time.Minute * 10)
	return edb
}

// CloseDb 关闭数据库连接
func CloseDb() error {
	return edb.CloseDb()
}
