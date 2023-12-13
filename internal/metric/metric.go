package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	agentMaxRetries     = 3
	agentInitialBackoff = 1 * time.Second
	agentMaxBackoff     = 5 * time.Second
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
		pollCount++
		url := fmt.Sprintf("http://%s", address)
		metrics, err := collectMetricsByMemStat(lastMemStats)
		if err != nil {
			logger.Logger().Error("failed to collect metrics by mem stat", zap.Error(err))
			break
		}

		err = sendUpdateMetricsRequest(url, metrics)
		if err != nil {
			logger.Logger().Error("failed to send gauge request to server", zap.Error(err))
			break
		}
	}
}

func pollLastMemStatByTicker(ticker *time.Ticker) {
	for t := range ticker.C {
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		lastMemStats = &memStats
		logger.Logger().Info("metric", zap.String("ticker", t.String()))
	}
}

func sendUpdateMetricsRequest(url string, metric []*storage.Metric) error {
	reqBodyBytes, err := json.Marshal(metric)
	if err != nil {
		logger.Logger().Error("failed to send metric: matshal request body", zap.Error(err))
		return err
	}
	compressedBytes, err := compressor.Compress(reqBodyBytes)
	if err != nil {
		logger.Logger().Error("failed to compress", zap.String("metric_json",
			string(reqBodyBytes)), zap.Error(err))
		return err
	}
	var request *http.Request
	request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/updates/", url),
		bytes.NewBuffer(compressedBytes),
	)
	if err != nil {
		logger.Logger().Error("failed to create new request",
			zap.String("metric_json", string(reqBodyBytes)), zap.Error(err))
		return err
	}
	request.Header.Set("Content-Type", handlers.UpdateMetricContentType)
	request.Header.Set("Accept-Encoding", compressor.AcceptEncoding)
	request.Header.Set("Content-Encoding", compressor.ContentEncoding)
	_ = retry.Do(
		func() error {
			var resp *http.Response
			resp, err = http.DefaultClient.Do(request)
			if err != nil {
				logger.Logger().Error("error to do request",
					zap.String("request.uri", request.RequestURI),
					zap.String("request.body", string(reqBodyBytes)),
					zap.Error(err),
				)
				return err
			}
			if resp == nil {
				return nil
			}
			if isRetriableHttpStatus(resp.StatusCode) {
				logger.Logger().Error("error response status",
					zap.String("request.uri", request.RequestURI),
					zap.String("request.body", string(reqBodyBytes)),
					zap.String("response.status", resp.Status),
					zap.Error(err),
				)
				return fmt.Errorf("response status %s code %d",
					resp.Status, resp.StatusCode)
			}

			var respBodyBytes []byte
			respBodyBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				logger.Logger().Error("error read response body",
					zap.String("request.uri", request.RequestURI),
					zap.String("request.body", string(reqBodyBytes)),
					zap.String("response.status", resp.Status),
					zap.Error(err),
				)
				// don't need to retry this error
				return nil
			}
			logger.Logger().Debug("read response body",
				zap.String("request.uri", request.RequestURI),
				zap.String("request.body", string(reqBodyBytes)),
				zap.String("response.status", resp.Status),
				zap.String("response.decompressed_body", string(respBodyBytes)),
			)
			err = resp.Body.Close()
			// don't need to retry this error
			return nil
		},
		retry.Attempts(agentMaxRetries),
		retry.Delay(agentInitialBackoff),
		retry.MaxDelay(agentMaxBackoff),
	)
	if err != nil {
		logger.Logger().Error(
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

	metrics = append(metrics, &storage.Metric{
		ID:    "PollCount",
		MType: storage.MTypeCounter,
		Delta: &pollCount,
	})
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

func isRetriableHttpStatus(statusCode int) bool {
	return statusCode == 500 && statusCode < 503
}
