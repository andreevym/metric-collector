package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/model"
	"github.com/go-chi/chi"
)

// GetMetricByTypeAndNameHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s Server) GetMetricByTypeAndNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		if metricType == model.MetricTypeCounter {
			v, err := counter.Get(s.CounterStorage(), metricName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			io.WriteString(w, v)
		} else if metricType == model.MetricTypeGauge {
			v, err := gauge.Get(s.GaugeStorage(), metricName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			io.WriteString(w, strings.Join(v, ", "))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
