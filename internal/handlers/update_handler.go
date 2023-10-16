package handlers

import (
	"net/http"

	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/go-chi/chi"
)

const (
	UpdateMetricContentType = "text/plain"
)

// UpdateMetricHandler method for insert or update metrics
// example request url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s ServiceHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdateMetricContentType)

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	err := multistorage.SaveMetric(s.storage, metricName, metricType, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
