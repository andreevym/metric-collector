package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/server"
)

func main() {
	cfg, err := config.ServerParse()
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("server config can't be nil")
	}

	server.Start(cfg.Address)
}
