package handlers

import (
	"encoding/json"
	"errors"
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
		logger.Logger().Error("error", zap.Error(err))
		return
	}

	var metrics []*storage.Metric
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		logger.Logger().Error("err", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metrics) == 0 {
		return
	}

	result := map[string]*storage.Metric{}
	for _, metric := range metrics {
		found, ok := result[metric.ID+metric.MType]
		if ok && found != nil && metric.MType == storage.MTypeCounter {
			newDelta := *metric.Delta + *found.Delta
			metric.Delta = &newDelta
		}

		result[metric.ID+metric.MType] = metric
	}

	metricsR := map[string]storage.MetricR{}
	for _, metric := range metrics {
		found, err := s.storage.Read(r.Context(), metric.ID, metric.MType)
		if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
			logger.Logger().Error("failed update metric",
				zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if found != nil && metric.MType == storage.MTypeCounter {
			newDelta := *metric.Delta + *found.Delta
			metric.Delta = &newDelta
		}

		metricsR[metric.ID+metric.MType] = storage.MetricR{
			Metric:   metric,
			IsExists: found != nil,
		}
	}

	err = s.storage.CreateAll(r.Context(), metricsR)
	if err != nil {
		logger.Logger().Error("err", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
	}
}
