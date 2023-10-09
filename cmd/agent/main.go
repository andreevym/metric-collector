package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/caarlos0/env"
)

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
	pollCount    int
	lastMemStats *runtime.MemStats
)

const (
	// ContentType Заголовок — Content-Type: text/plain.
	ContentType = "text/plain"
)

func main() {
	parseFlags()

	var config EnvConfig
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Address != "" {
		flagRunAddr = config.Address
	}

	// Обновлять метрики из пакета runtime с заданной частотой: pollInterval — 2 секунды.
	if config.PollInterval != 0 {
		flagPollInterval = config.PollInterval
	}
	tickerMetric := time.NewTicker(flagPollInterval * time.Second)

	// Отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
	if config.ReportInterval != 0 {
		flagReportInterval = config.ReportInterval
	}
	tickerReport := time.NewTicker(flagReportInterval * time.Second)

	go collectRuntimeMetric(tickerMetric)
	go sendMetricToServer(tickerReport, flagRunAddr)

	// время жизни клиента для сбора метрик
	time.Sleep(time.Minute)
}

func sendMetricToServer(tickerReport *time.Ticker, metricServerURL string) {
	for a := range tickerReport.C {
		randomValue := rand.Int()
		fmt.Printf("- report randomValue: %v", randomValue)
		fmt.Printf("- report pollCount %v\n", pollCount)
		fmt.Printf("- report lastMemStats %v\n", lastMemStats)
		fmt.Printf("- report %s\n", a.String())
		resp, err := http.Post(metricServerURL, ContentType, nil)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to send metric: invalid send http post: %v", err)
			break
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("failed to send metri:, invalid response status. found status %s", resp.Status)
			break
		}
	}
}

func collectRuntimeMetric(tickerMetric *time.Ticker) {
	for a := range tickerMetric.C {
		pollCount++
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		lastMemStats = &memStats
		fmt.Printf("+ metric %s\n", a.String())
	}
}
