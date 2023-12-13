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

	metric, err := buildMetricByHttpRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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

	bytes, err := json.Marshal(&metric)
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

func buildMetricByHttpRequest(r *http.Request) (*storage.Metric, error) {
	metric, err := buildMetricByBody(r.Body)
	if err == nil && metric != nil {
		return metric, nil
	}

	return buildMetricByParam(r)
}

func buildMetricByBody(body io.ReadCloser) (*storage.Metric, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		logger.Logger().Error("error", zap.Error(err))
	}

	if len(bytes) == 0 {
		return nil, errors.New("body len is empty")
	}

	metric := &storage.Metric{}
	err = json.Unmarshal(bytes, &metric)
	return metric, err
}

func buildMetricByParam(r *http.Request) (*storage.Metric, error) {
	metric := &storage.Metric{}
	metric.MType = chi.URLParam(r, "metricType")
	metric.ID = chi.URLParam(r, "metricName")
	v := chi.URLParam(r, "metricValue")
	if metric.MType == storage.MTypeCounter {
		delta, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		metric.Delta = &delta
	} else if metric.MType == storage.MTypeGauge {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		metric.Value = &value
	} else {
		return nil, errors.New("unknown type")
	}

	return metric, nil
}
