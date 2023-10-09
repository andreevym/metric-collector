package main

import "github.com/andreevym/metric-collector/internal/server"

func main() {
	parseFlags()

	server.StartServer(flagRunAddr)
}
