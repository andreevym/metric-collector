package handlers

import (
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/model"
	"github.com/andreevym/metric-collector/internal/repository"
)

const (
	idxMetricType       = 2
	idxMetricName       = 3
	idxMetricValue      = 4
	lenArgsUpdateMethod = 5
)

// UpdateHandler method for insert or UpdateHandler metrics
// example url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandler(counterStorage repository.Storage, gaugeStorage repository.Storage) http.HandlerFunc {
	if counterStorage == nil {
		panic("counter storage can't be nil")
	}
	if gaugeStorage == nil {
		panic("gauge storage can't be nil")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		values := strings.Split(r.URL.Path, "/")
		if len(values) != lenArgsUpdateMethod {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metricName := values[idxMetricName]
		metricValue := values[idxMetricValue]
		metricType := values[idxMetricType]
		if metricName == "" ||
			metricValue == "" ||
			metricType == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if metricType == model.MetricTypeCounter {
			err := counter.Validate(metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = counter.Store(counterStorage, metricName, metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if metricType == model.MetricTypeGauge {
			err := gauge.Validate(metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = gauge.Store(gaugeStorage, metricName, metricValue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
