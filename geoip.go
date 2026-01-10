package netguard

import (
	"net/netip"
	"sync"

	"github.com/iotames/netguard/log"
	"github.com/oschwald/maxminddb-golang/v2"
	// "github.com/oschwald/geoip2-golang/v2"
)

var (
	geoipdbFile = "GeoLite2-City.mmdb"
	geoipDb     *maxminddb.Reader
	once        sync.Once
)

func SetGeoipDb(file string) error {
	var err error
	geoipdbFile = file
	geoipDb, err = maxminddb.Open(geoipdbFile)
	return err
}

func getGeoipDb() *maxminddb.Reader {
	var err error
	once.Do(func() {
		if geoipDb == nil {
			geoipDb, err = maxminddb.Open(geoipdbFile)
			if err != nil {
				log.Error("error", err.Error())
				panic(err)
			}
		}
	})
	return geoipDb
}

type GeoIpInfo struct {
	CountryCode string
	Country     string
	City        string
}

func GetIpGeo(remoteIP string) GeoIpInfo {
	// https://github.com/P3TERX/GeoLite.mmdb/releases
	// https://github.com/P3TERX/GeoLite.mmdb/releases/download/2025.09.22/GeoLite2-City.mmdb
	// db, err := maxminddb.Open("GeoLite2-City.mmdb")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	db := getGeoipDb()

	ip, err := netip.ParseAddr(remoteIP)
	if err != nil {
		log.Warn("parse remoteIP fail", "remoteIP", remoteIP, "error", err.Error())
		return GeoIpInfo{}
	}

	var record struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		Subdivisions []struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"subdivisions"`
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
	}

	err = db.Lookup(ip).Decode(&record)
	if err != nil {
		log.Warn("geoip lookup fail", "remoteIP", remoteIP, "error", err.Error())
		return GeoIpInfo{}
	}

	// fmt.Printf("Country: %s (%s)\n", record.Country.Names["zh-CN"], record.Country.ISOCode)
	// fmt.Printf("City: %s\n", record.City.Names["zh-CN"])
	// if len(record.Subdivisions) > 0 {
	// 	fmt.Printf("Subdivision: %s\n", record.Subdivisions[0].Names["zh-CN"])
	// }
	// fmt.Printf("CountryInfo: %+v\n", record.Country.Names)
	// fmt.Printf("CityInfo: %+v\n", record.City.Names)
	return GeoIpInfo{record.Country.ISOCode, record.Country.Names["zh-CN"], record.City.Names["zh-CN"]}
}
