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
		logger.Log.Error("error", zap.Error(err))
	}

	m := &storage.Metric{}
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		m.MType = chi.URLParam(r, "metricType")
		m.ID = chi.URLParam(r, "metricName")
		if m.MType == storage.MTypeCounter {
			v := chi.URLParam(r, "metricValue")
			delta, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Delta = &delta
		} else if m.MType == storage.MTypeGauge {
			v := chi.URLParam(r, "metricValue")
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			m.Value = &value
		}
	}

	if m.MType != storage.MTypeGauge && m.MType != storage.MTypeCounter {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundValue, err := s.storage.Read(r.Context(), m.ID)
	if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
		logger.Log.Error("failed update metric",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if foundValue == nil {
		err = s.storage.Create(r.Context(), m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		if m.MType == storage.MTypeCounter {
			newDelta := *m.Delta + *foundValue.Delta
			m.Delta = &newDelta
		}
		err = s.storage.Update(r.Context(), m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	bytes, err = json.Marshal(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
