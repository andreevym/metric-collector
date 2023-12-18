package metricagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/hash"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/utils"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
)

const retryAttempts = 100000

var (
	// PollCount (тип counter) — счётчик, увеличивающийся на 1
	// при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
	pollCount    int64
	lastMemStats *runtime.MemStats
)

// sendLastMemStats send metric to server by ticker and address
func sendLastMemStats(ctx context.Context, secretKey string, ticker *time.Ticker, address string) {
	for t := range ticker.C {
		logger.Logger().Info("sendLastMemStats",
			zap.String("ticker", t.String()),
			zap.Int64("pollCount", pollCount),
		)
		pollCount++
		metrics, err := collectMetricsByMemStat(lastMemStats, pollCount)
		if err != nil {
			logger.Logger().Error("failed to collect metrics by mem stat", zap.Error(err))
			break
		}

		err = sendUpdateMetricsRequest(ctx, secretKey, address, metrics)
		if err != nil {
			logger.Logger().Error("failed to send update request with last metric", zap.Error(err))
			break
		}
	}
}

func sendUpdateMetricsRequest(ctx context.Context, secretKey string, address string, metric []*storage.Metric) error {
	b, err := json.Marshal(metric)
	if err != nil {
		logger.Logger().Error("failed to marshal request body", zap.Error(err))
		return err
	}
	compressedBytes, err := compressor.Compress(b)
	if err != nil {
		logger.Logger().Error("failed to compress",
			zap.String("request body", string(b)), zap.Error(err))
		return err
	}
	_ = retry.Do(
		func() error {
			var request *http.Request
			request, err = http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				fmt.Sprintf("http://%s/updates/", address),
				bytes.NewBuffer(compressedBytes),
			)
			if err != nil {
				logger.Logger().Error("failed to create new request",
					zap.String("request body", string(b)), zap.Error(err))
				return err
			}
			request.Header.Set("Content-Type", handlers.UpdateMetricContentType)
			request.Header.Set("Accept-Encoding", compressor.AcceptEncoding)
			request.Header.Set("Content-Encoding", compressor.ContentEncoding)
			if len(secretKey) != 0 {
				request.Header.Set("HashSHA256", hash.Hash(compressedBytes, secretKey))
			}
			var resp *http.Response
			resp, err = http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			if resp == nil {
				return nil
			}
			if isRetriableStatus(resp.StatusCode) {
				logger.Logger().Error("error response status",
					zap.String("request.URL", request.URL.String()),
					zap.Error(err),
				)
				return fmt.Errorf("response status %s code %d",
					resp.Status, resp.StatusCode)
			}

			var respBodyBytes []byte
			respBodyBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				logger.Logger().Error("error read response body",
					zap.String("request.URL", request.URL.String()),
					zap.String("request.body", string(b)),
					zap.String("response.status", resp.Status),
					zap.Error(err),
				)
				// don't need to retry this error
				return nil
			}
			logger.Logger().Debug("read response body",
				zap.String("request.URL", request.URL.String()),
				zap.String("request.body", string(b)),
				zap.String("response.status", resp.Status),
				zap.String("response.decompressed_body", string(respBodyBytes)),
			)
			err = resp.Body.Close()
			// don't need to retry this error
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error to send request",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	if err != nil {
		logger.Logger().Error(
			"send request error",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func isRetriableStatus(statusCode int) bool {
	return statusCode == 500 && statusCode < 503
}
