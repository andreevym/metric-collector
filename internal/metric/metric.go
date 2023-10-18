package metric

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andreevym/metric-collector/internal/handlers"
)

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
	pollCount    int
	lastMemStats *runtime.MemStats
)

func Start(pollDuration time.Duration, reportDuration time.Duration, address string) {
	tickerPoll := time.NewTicker(pollDuration)
	tickerReport := time.NewTicker(reportDuration)

	go pollLastMemStatByTicker(tickerPoll)
	go sendByTickerAndAddress(tickerReport, address)

	// время жизни клиента для сбора метрик
	time.Sleep(time.Minute)
}

// sendByTickerAndAddress send metric to server by ticker and address
func sendByTickerAndAddress(ticker *time.Ticker, address string) {
	for a := range ticker.C {
		randomValue := rand.Int()
		fmt.Printf("- report randomValue: %v", randomValue)
		fmt.Printf("- report pollCount %v\n", pollCount)
		fmt.Printf("- report lastMemStats %v\n", lastMemStats)
		fmt.Printf("- report %s\n", a.String())
		url := fmt.Sprintf("http://%s", address)
		resp, err := http.Post(url, handlers.UpdateMetricContentType, nil)
		if err != nil {
			fmt.Printf("failed to send metric: invalid send http post: %v", err)
			break
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to handle response from server: close resp body: %v", err)
			break
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("invalid response status. found status code %d but expected %d",
				resp.StatusCode, http.StatusOK)
			break
		}
	}
}

func pollLastMemStatByTicker(ticker *time.Ticker) {
	for a := range ticker.C {
		pollCount++
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		lastMemStats = &memStats
		fmt.Printf("+ metric %s\n", a.String())
	}
}
