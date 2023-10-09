package handlers

import (
	"net/http"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/model"
	"github.com/go-chi/chi"
)

// UpdateMetricHandler method for insert or update metrics
// example request url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s Server) UpdateMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		switch metricType {
		case model.MetricTypeCounter:
			err := counter.Validate(metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = counter.Store(s.CounterStorage(), metricName, metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case model.MetricTypeGauge:
			err := gauge.Validate(metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = gauge.Store(s.GaugeStorage(), metricName, metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
