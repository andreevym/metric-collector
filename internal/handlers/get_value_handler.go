package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetValueHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	v, err := s.metricStorage.GetMetric(metricType, metricName)
	if err != nil || v == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		_, err := io.WriteString(w, v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// GetPingHandler ping database
func (s ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	err := s.dbClient.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
