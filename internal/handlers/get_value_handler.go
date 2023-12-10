package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

// GetValueHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	if metricType != storage.MTypeGauge && metricType != storage.MTypeCounter {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "metricName")

	v, err := s.storage.Read(r.Context(), id, metricType)
	if err != nil || v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var res string
	switch v.MType {
	case storage.MTypeCounter:
		if v.Delta == nil {
			logger.Log.Error("delta can't be nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res = strconv.FormatInt(*v.Delta, 10)
		_, err = io.WriteString(w, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		return
	case storage.MTypeGauge:
		if v.Value == nil {
			logger.Log.Error("value can't be nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strconv.FormatFloat(*v.Value, 'f', -1, 64)
		_, err = io.WriteString(w, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetPingHandler ping database
func (s ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	err := s.dbClient.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
