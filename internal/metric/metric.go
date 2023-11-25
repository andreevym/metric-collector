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
		collectedMetrics, err := collectMetricsByMemStat(lastMemStats)
		if err != nil {
			logger.Log.Error("failed to collect metrics by mem stat", zap.Error(err))
			break
		}
		err = sendPollCount(url)
		if err != nil {
			logger.Log.Error("failed to send counter request to server", zap.Error(err))
			break
		}
		for _, m := range collectedMetrics {
			err = sendGauge(url, m)
			if err != nil {
				logger.Log.Error("failed to send gauge request to server", zap.Error(err))
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
		logger.Log.Info("+ metric\n", zap.Int64("pollCount", pollCount), zap.String("ticker", a.String()))
	}
}

func sendGauge(url string, m *storage.Metric) error {
	m.Delta = nil
	return sendUpdateMetricsRequest(url, m)
}

func sendPollCount(url string) error {
	m := &storage.Metric{
		ID:    PollCount,
		MType: storage.MTypeCounter,
		Delta: &pollCount,
		Value: nil,
	}
	return sendUpdateMetricsRequest(url, m)
}

func sendUpdateMetricsRequest(url string, metric *storage.Metric) error {
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
				fmt.Sprintf("%s/update", url),
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
				logger.Log.Error(
					"send request error",
					zap.String("url", url),
					zap.Error(err),
				)
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
			"retry error",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func collectMetricsByMemStat(stats *runtime.MemStats) ([]*storage.Metric, error) {
	metrics := make([]*storage.Metric, 0)

	MustAppendGaugeMetricFloat64(metrics, RandomValue, rand.Float64())
	MustAppendGaugeMetricUint64(metrics, Alloc, stats.Alloc)
	MustAppendGaugeMetricUint64(metrics, BuckHashSys, stats.BuckHashSys)
	MustAppendGaugeMetricUint64(metrics, Frees, stats.Frees)
	MustAppendGaugeMetricFloat64(metrics, GCCPUFraction, stats.GCCPUFraction)
	MustAppendGaugeMetricUint64(metrics, GCSys, stats.GCSys)
	MustAppendGaugeMetricUint64(metrics, HeapAlloc, stats.HeapAlloc)
	MustAppendGaugeMetricUint64(metrics, HeapIdle, stats.HeapIdle)
	MustAppendGaugeMetricUint64(metrics, HeapInuse, stats.HeapInuse)
	MustAppendGaugeMetricUint64(metrics, HeapObjects, stats.HeapObjects)
	MustAppendGaugeMetricUint64(metrics, HeapReleased, stats.HeapReleased)
	MustAppendGaugeMetricUint64(metrics, HeapSys, stats.HeapSys)
	MustAppendGaugeMetricUint64(metrics, LastGC, stats.LastGC)
	MustAppendGaugeMetricUint64(metrics, Lookups, stats.Lookups)
	MustAppendGaugeMetricUint64(metrics, MCacheInuse, stats.MCacheInuse)
	MustAppendGaugeMetricUint64(metrics, MCacheSys, stats.MCacheSys)
	MustAppendGaugeMetricUint64(metrics, MSpanInuse, stats.MSpanInuse)
	MustAppendGaugeMetricUint64(metrics, MSpanSys, stats.MSpanSys)
	MustAppendGaugeMetricUint64(metrics, Mallocs, stats.Mallocs)
	MustAppendGaugeMetricUint64(metrics, NextGC, stats.NextGC)
	MustAppendGaugeMetricUint32(metrics, NumForcedGC, stats.NumForcedGC)
	MustAppendGaugeMetricUint32(metrics, NumGC, stats.NumGC)
	MustAppendGaugeMetricUint64(metrics, OtherSys, stats.OtherSys)
	MustAppendGaugeMetricUint64(metrics, PauseTotalNs, stats.PauseTotalNs)
	MustAppendGaugeMetricUint64(metrics, StackInuse, stats.StackInuse)
	MustAppendGaugeMetricUint64(metrics, StackSys, stats.StackSys)
	MustAppendGaugeMetricUint64(metrics, Sys, stats.Sys)
	MustAppendGaugeMetricUint64(metrics, TotalAlloc, stats.TotalAlloc)
	return metrics, nil
}

func MustAppendGaugeMetricUint64(metrics []*storage.Metric, id string, v uint64) {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
}

func MustAppendGaugeMetricUint32(metrics []*storage.Metric, id string, v uint32) {
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 32)
	if err != nil {
		panic(err)
	}
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
}

func MustAppendGaugeMetricFloat64(metrics []*storage.Metric, id string, f float64) {
	metrics = append(metrics, &storage.Metric{
		ID:    id,
		MType: storage.MTypeGauge,
		Value: &f,
	})
}
