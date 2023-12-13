package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	UpdateMetricContentType = "application/json"
	ValueMetricContentType  = "application/json"
)

// PostUpdateHandler method for insert or update metrics
// example request url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s ServiceHandlers) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdateMetricContentType)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger().Error("error", zap.Error(err))
	}

	metric := &storage.Metric{}
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		metric.MType = chi.URLParam(r, "metricType")
		metric.ID = chi.URLParam(r, "metricName")
		v := chi.URLParam(r, "metricValue")
		if metric.MType == storage.MTypeCounter {
			delta, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metric.Delta = &delta
		} else if metric.MType == storage.MTypeGauge {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metric.Value = &value
		}
	}

	if metric.MType != storage.MTypeGauge && metric.MType != storage.MTypeCounter {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundValue, err := s.storage.Read(r.Context(), metric.ID, metric.MType)
	if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
		logger.Logger().Error("failed update metric",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if foundValue == nil {
		err = s.storage.Create(r.Context(), metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		if metric.MType == storage.MTypeCounter {
			newDelta := *metric.Delta + *foundValue.Delta
			metric.Delta = &newDelta
		}
		err = s.storage.Update(r.Context(), metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	bytes, err = json.Marshal(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
