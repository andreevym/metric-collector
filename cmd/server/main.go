package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/server"
)

func main() {
	config.ServerFlags()

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
