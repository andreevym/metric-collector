package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/multistorage"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gauge := map[string]string{}
	counter := map[string]string{}
	for _, metric := range metrics {
		if metric.MType == multistorage.MetricTypeGauge {
			if metric.Value == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metricValue := fmt.Sprintf("%v", *metric.Value)

			gauge[metric.ID] = metricValue
		} else if metric.MType == multistorage.MetricTypeCounter {
			if metric.Delta == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metricValue := strconv.FormatInt(*metric.Delta, 10)

			counter[metric.ID] = metricValue
		}
	}

	if len(gauge) > 0 {
		err = s.metricStorage.SaveMetrics(multistorage.MetricTypeGauge, gauge)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if len(counter) > 0 {
		err = s.metricStorage.SaveMetrics(multistorage.MetricTypeCounter, counter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
