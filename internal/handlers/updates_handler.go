package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"go.uber.org/zap"
)

const (
	UpdatesMetricContentType = "application/json"
)

func (s ServiceHandlers) PostUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdatesMetricContentType)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("error", zap.Error(err))
		return
	}

	var metrics []*storage.Metric
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricsR := map[string]*storage.MetricR{}
	for _, m := range metrics {
		foundMetric, err := s.storage.Read(r.Context(), m.ID)
		if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
			logger.Log.Error("failed update metric",
				zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if m.MType == storage.MTypeGauge {
			if m.Value == nil {
				logger.Log.Error("failed update metric",
					zap.Error(fmt.Errorf("value can't be nil for id %s", m.ID)))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			metricsR[m.ID] = &storage.MetricR{
				Metric:   m,
				IsExists: false,
			}
		} else if m.MType == storage.MTypeCounter {
			if m.Delta == nil {
				logger.Log.Error("err", zap.Error(fmt.Errorf("delta can't be nil for id %s", m.ID)))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var existsMetricValue *storage.Metric

			cachedMetric, ok := metricsR[m.ID]
			if ok {
				existsMetricValue = cachedMetric.Metric
			} else if foundMetric == nil {
				metricsR[m.ID] = &storage.MetricR{
					Metric:   m,
					IsExists: false,
				}
				break
			} else if foundMetric != nil {
				existsMetricValue = foundMetric
			}

			if existsMetricValue != nil {
				newVal := *existsMetricValue.Delta + *m.Delta
				m.Delta = &newVal
			}

			metricsR[m.ID] = &storage.MetricR{
				Metric:   m,
				IsExists: true,
			}
		}
	}

	if len(metricsR) > 0 {
		err = s.storage.CreateAll(r.Context(), metricsR)
		if err != nil {
			logger.Log.Error("err", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
