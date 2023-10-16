package handlers

import (
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/go-chi/chi"
)

// GetMetricByTypeAndNameHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) GetMetricByTypeAndNameHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	v, err := multistorage.GetMetric(s.storage, metricType, metricName)
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
