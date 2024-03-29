package metricagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

const retryAttempts = 3

// sendLastMemStats send metric to server by ticker and address
func sendMetric(
	ctx context.Context,
	inputCh chan []*storage.Metric,
	secretKey string,
	reportDuration time.Duration,
	address string,
) {
	ticker := time.NewTicker(reportDuration)
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		case metrics, ok := <-inputCh:
			if ok {
				err := sendRequest(ctx, secretKey, address, metrics)
				if err != nil {
					logger.Logger().Error("failed to send request", zap.Error(err))
					break
				}
			}
		default:
		}
	}
}

func sendRequest(
	ctx context.Context,
	secretKey string,
	address string,
	metric []*storage.Metric,
) error {
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
			request, err = http.NewRequest(
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
				request.Header.Set("HashSHA256", hash.EncodeHash(compressedBytes, secretKey))
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
		retry.Context(ctx),
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
