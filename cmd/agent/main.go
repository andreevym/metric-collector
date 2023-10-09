package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
	pollCount    int
	lastMemStats *runtime.MemStats
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second

	// Формат данных — http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>.
	// Адрес сервера — http://localhost:8080.
	metricServerURL = "http://localhost:8080"
	// ContentType Заголовок — Content-Type: text/plain.
	ContentType = "text/plain"
)

func main() {
	// Обновлять метрики из пакета runtime с заданной частотой: pollInterval — 2 секунды.
	tickerMetric := time.NewTicker(pollInterval)
	// Отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
	tickerReport := time.NewTicker(reportInterval)

	go collectRuntimeMetric(tickerMetric)
	go sendMetricToServer(tickerReport)

	// время жизни клиента для сбора метрик
	time.Sleep(time.Minute)
}

func sendMetricToServer(tickerReport *time.Ticker) {
	for a := range tickerReport.C {
		randomValue := rand.Int()
		fmt.Printf("- report randomValue: %v", randomValue)
		fmt.Printf("- report pollCount %v\n", pollCount)
		fmt.Printf("- report lastMemStats %v\n", lastMemStats)
		fmt.Printf("- report %s\n", a.String())
		resp, err := http.Post(metricServerURL, ContentType, nil)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic("resp.StatusCode != http.StatusOK")
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
