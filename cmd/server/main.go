package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/server"
)

func init() {
	config.ServerFlags()
}

func main() {
	cfg, err := config.ServerParse()
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("server config can't be nil")
	}

	err = server.Start(cfg.Address)
	if err != nil {
		log.Fatal(err)
	}
}
