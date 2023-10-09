package handlers

import (
	"io"
	"net/http"

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

		switch metricType {
		case model.MetricTypeCounter:
			v, err := counter.Get(s.CounterStorage(), metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v != "" {
				io.WriteString(w, v)
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case model.MetricTypeGauge:
			v, err := gauge.Get(s.GaugeStorage(), metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if len(v) != 0 {
				io.WriteString(w, v[len(v)-1])
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
