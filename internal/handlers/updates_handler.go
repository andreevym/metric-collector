package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/multistorage"
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

	var metrics []Metrics
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gauge := map[string]*storage.Metric{}
	counter := map[string]*storage.Metric{}
	for _, metric := range metrics {
		if metric.MType == multistorage.MetricTypeGauge {
			if metric.Value == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metricValue := fmt.Sprintf("%v", *metric.Value)

			gauge[metric.ID] = &storage.Metric{
				Value:    metricValue,
				IsExists: false,
			}
		} else if metric.MType == multistorage.MetricTypeCounter {
			if metric.Delta == nil {
				logger.Log.Error("err", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			isExists := true
			addedMetricValue := strconv.FormatInt(*metric.Delta, 10)

			var existsMetricValue string

			cachedMetric, ok := counter[metric.ID]
			if ok {
				existsMetricValue = cachedMetric.Value
			} else {
				existsMetricValue, err = s.metricStorage.GetMetric(multistorage.MetricTypeCounter, metric.ID)
				if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
					logger.Log.Error("err", zap.Error(err))
					w.WriteHeader(http.StatusBadRequest)
					return
				} else if err != nil && errors.Is(err, storage.ErrValueNotFound) {
					existsMetricValue = "0"
					isExists = false
				}
			}

			existsMetricVal, err := strconv.ParseInt(existsMetricValue, 10, 64)
			if err != nil {
				logger.Log.Error("err", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			v, err := strconv.ParseInt(addedMetricValue, 10, 64)
			if err != nil {
				logger.Log.Error("err", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			newVal := strconv.FormatInt(existsMetricVal+v, 10)

			counter[metric.ID] = &storage.Metric{
				Value:    newVal,
				IsExists: isExists,
			}
		}
	}

	if len(gauge) > 0 {
		err = s.metricStorage.SaveMetrics(multistorage.MetricTypeGauge, gauge)
		if err != nil {
			logger.Log.Error("err", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if len(counter) > 0 {
		err = s.metricStorage.SaveMetrics(multistorage.MetricTypeCounter, counter)
		if err != nil {
			logger.Log.Error("err", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
