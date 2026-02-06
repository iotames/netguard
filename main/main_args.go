package main

import (
	"flag"

	"github.com/iotames/netguard/conf"
)

var Devname string
var ListDev, V, VersionV bool
var Port int

func parseArgs() {
	flag.StringVar(&Devname, "devname", "", `netguard.exe --devname="\Device\NPF_{3757BF1E-96B9-441B-8D4B-95EAB49ECA36}"`)
	flag.BoolVar(&ListDev, "listdev", false, "netguard.exe --listdev")
	flag.IntVar(&Port, "port", conf.WebServerPort, "netguard.exe --port=8080")
	flag.BoolVar(&V, "v", false, "netguard.exe --v")
	flag.BoolVar(&VersionV, "version", false, "netguard.exe --version")
	flag.Parse()
}
