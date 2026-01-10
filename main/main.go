package main

import (
	"log/slog"

	"github.com/iotames/netguard"
	"github.com/iotames/netguard/log"
)

func main() {
	log.SetLevel(slog.LevelInfo)
	// f, err := log.SetLogWriterByFile("netguard.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	netguard.Run()
}
