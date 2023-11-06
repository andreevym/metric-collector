package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
)

const (
	defaultRetryCount = 30
	defaultRetryWait  = 50 * time.Millisecond
)

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
	pollCount    int64
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
	for range ticker.C {
		url := fmt.Sprintf("http://%s", address)
		stats, err := collectMetricsByMemStat(lastMemStats)
		if err != nil {
			fmt.Printf("failed to collect metrics by mem stat: %v", err)
			break
		}
		err = sendPollCount(url)
		if err != nil {
			fmt.Printf("failed to send counter request to server: %v", err)
			break
		}
		for _, metrics := range stats {
			err = sendGauge(metrics, url)
			if err != nil {
				fmt.Printf("failed to send gauge request to server: %v", err)
				break
			}
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

func sendGauge(metrics handlers.Metrics, url string) error {
	metrics.Delta = nil
	return sendUpdateMetricsRequest(url, metrics)
}

func sendPollCount(url string) error {
	metrics := handlers.Metrics{
		ID:    PollCount,
		MType: multistorage.MetricTypeCounter,
		Delta: &pollCount,
		Value: nil,
	}
	return sendUpdateMetricsRequest(url, metrics)
}

func sendUpdateMetricsRequest(url string, metrics handlers.Metrics) error {
	b, err := json.Marshal(metrics)
	if err != nil {
		fmt.Printf("failed to send metric: matshal request body: %v", err)
		return err
	}
	var resp *http.Response
	err = retry.Do(
		func() error {
			var err error
			resp, err = http.Post(fmt.Sprintf("%s/update", url), handlers.UpdateMetricContentType, bytes.NewBuffer(b))
			return err
		},
		retry.Attempts(defaultRetryCount),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		zap.Error(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func collectMetricsByMemStat(stats *runtime.MemStats) ([]handlers.Metrics, error) {
	m := make([]handlers.Metrics, 0)

	f := rand.Float64()
	m = append(m, handlers.Metrics{
		ID:    RandomValue,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err := strconv.ParseFloat(fmt.Sprintf("%v", stats.Alloc), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    Alloc,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.BuckHashSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    BuckHashSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.Frees), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    Frees,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	m = append(m, handlers.Metrics{
		ID:    GCCPUFraction,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &stats.GCCPUFraction,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.GCSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    GCSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &stats.GCCPUFraction,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapAlloc), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapAlloc,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapIdle), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapIdle,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapInuse), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapInuse,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapObjects), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapObjects,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapReleased), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapReleased,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.HeapSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    HeapSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.LastGC), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    LastGC,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.Lookups), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    Lookups,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.MCacheInuse), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    MCacheInuse,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.MCacheSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    MCacheSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.MSpanInuse), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    MSpanInuse,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.MSpanSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    MSpanSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.Mallocs), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    Mallocs,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.NextGC), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    NextGC,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.NumForcedGC), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    NumForcedGC,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.NumGC), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    NumGC,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.OtherSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    OtherSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.PauseTotalNs), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    PauseTotalNs,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.StackInuse), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    StackInuse,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.StackSys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    StackSys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.Sys), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    Sys,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	f, err = strconv.ParseFloat(fmt.Sprintf("%v", stats.TotalAlloc), 64)
	if err != nil {
		return nil, err
	}
	m = append(m, handlers.Metrics{
		ID:    TotalAlloc,
		MType: multistorage.MetricTypeGauge,
		Delta: nil,
		Value: &f,
	})

	return m, nil
}
