package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
)

const (
	defaultRetryCount = 100
)

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
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
		collectedMetrics, err := collectMetricsByMemStat(lastMemStats)
		if err != nil {
			logger.Log.Error("failed to collect metrics by mem stat", zap.Error(err))
			break
		}

		i := int64(1)
		metricPollCount := &storage.Metric{
			ID:    "PollCount",
			MType: storage.MTypeCounter,
			Delta: &i,
			Value: nil,
		}
		collectedMetrics = append(collectedMetrics, metricPollCount)
		err = sendUpdateMetricsRequest(url, collectedMetrics)
		if err != nil {
			logger.Log.Error("failed to send gauge request to server", zap.Error(err))
			break
		}
	}
}

func pollLastMemStatByTicker(ticker *time.Ticker) {
	for a := range ticker.C {
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		lastMemStats = &memStats
		logger.Log.Info("+ metric\n", zap.String("ticker", a.String()))
	}
}

func sendUpdateMetricsRequest(url string, metric []*storage.Metric) error {
	b, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Error("failed to send metric: matshal request body", zap.Error(err))
		return err
	}
	err = retry.Do(
		func() error {
			compressedBytes, err := compressor.Compress(b)
			if err != nil {
				logger.Log.Error("failed to compress", zap.String("metric_json", string(b)), zap.Error(err))
				return err
			}
			request, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/updates/", url),
				bytes.NewBuffer(compressedBytes),
			)
			if err != nil {
				logger.Log.Error("failed to create new request", zap.String("metric_json", string(b)), zap.Error(err))
				return err
			}
			request.Header.Set("Content-Type", handlers.UpdateMetricContentType)
			request.Header.Set("Accept-Encoding", compressor.AcceptEncoding)
			request.Header.Set("Content-Encoding", compressor.ContentEncoding)
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			if resp != nil {
				bytes, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Log.Error("error read response body",
						zap.String("request.uri", request.RequestURI),
						zap.String("request.body", string(b)),
					)
					return err
				}
				logger.Log.Debug("read response body",
					zap.String("request.uri", request.RequestURI),
					zap.String("request.body", string(b)),
					zap.String("response.status", resp.Status),
					zap.String("response.decompressed_body", string(bytes)),
				)
				err = resp.Body.Close()
				if err != nil {
					return err
				}
			}
			return nil
		},
		retry.Attempts(defaultRetryCount),
		retry.Delay(500*time.Millisecond),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		logger.Log.Error(
			"send request error",
			zap.String("url", url),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func collectMetricsByMemStat(stats *runtime.MemStats) ([]*storage.Metric, error) {
	metrics := make([]*storage.Metric, 0)
	metrics = MustAppendGaugeMetricFloat64(metrics, "RandomValue", rand.Float64())
	metrics = MustAppendGaugeMetricUint64(metrics, "Alloc", stats.Alloc)
	metrics = MustAppendGaugeMetricUint64(metrics, "BuckHashSys", stats.BuckHashSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "Frees", stats.Frees)
	metrics = MustAppendGaugeMetricFloat64(metrics, "GCCPUFraction", stats.GCCPUFraction)
	metrics = MustAppendGaugeMetricUint64(metrics, "GCSys", stats.GCSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapAlloc", stats.HeapAlloc)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapIdle", stats.HeapIdle)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapInuse", stats.HeapInuse)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapObjects", stats.HeapObjects)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapReleased", stats.HeapReleased)
	metrics = MustAppendGaugeMetricUint64(metrics, "HeapSys", stats.HeapSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "LastGC", stats.LastGC)
	metrics = MustAppendGaugeMetricUint64(metrics, "Lookups", stats.Lookups)
	metrics = MustAppendGaugeMetricUint64(metrics, "MCacheInuse", stats.MCacheInuse)
	metrics = MustAppendGaugeMetricUint64(metrics, "MCacheSys", stats.MCacheSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "MSpanInuse", stats.MSpanInuse)
	metrics = MustAppendGaugeMetricUint64(metrics, "MSpanSys", stats.MSpanSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "Mallocs", stats.Mallocs)
	metrics = MustAppendGaugeMetricUint64(metrics, "NextGC", stats.NextGC)
	metrics = MustAppendGaugeMetricUint32(metrics, "NumForcedGC", stats.NumForcedGC)
	metrics = MustAppendGaugeMetricUint32(metrics, "NumGC", stats.NumGC)
	metrics = MustAppendGaugeMetricUint64(metrics, "OtherSys", stats.OtherSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "PauseTotalNs", stats.PauseTotalNs)
	metrics = MustAppendGaugeMetricUint64(metrics, "StackInuse", stats.StackInuse)
	metrics = MustAppendGaugeMetricUint64(metrics, "StackSys", stats.StackSys)
	metrics = MustAppendGaugeMetricUint64(metrics, "Sys", stats.Sys)
	metrics = MustAppendGaugeMetricUint64(metrics, "TotalAlloc", stats.TotalAlloc)
	return metrics, nil
}

func MustAppendGaugeMetricUint64(metrics []*storage.Metric, id string, v uint64) []*storage.Metric {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}

func MustAppendGaugeMetricUint32(metrics []*storage.Metric, id string, v uint32) []*storage.Metric {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 32)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}

func MustAppendGaugeMetricFloat64(metrics []*storage.Metric, id string, f float64) []*storage.Metric {
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
	return metrics
}
