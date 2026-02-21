package db

import (
	"github.com/iotames/easydb"
)

var edb *easydb.EasyDb

func SetDb(d *easydb.EasyDb) {
	edb = d
}
func GetDb() *easydb.EasyDb {
	return edb
}
