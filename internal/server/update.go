package server

import (
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/counter"
	"github.com/andreevym/metric-collector/internal/gauge"
	"github.com/andreevym/metric-collector/internal/model"
)

const (
	idxMetricType       = 2
	idxMetricName       = 3
	idxMetricValue      = 4
	lenArgsUpdateMethod = 5
)

// update method for insert or update metrics
// example url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s Server) update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	values := strings.Split(r.URL.Path, "/")
	if len(values) != lenArgsUpdateMethod {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricName := values[idxMetricName]
	metricValue := values[idxMetricValue]
	metricType := values[idxMetricType]

	if metricType == model.MetricTypeCounter {
		err := counter.StoreCounter(s.counterMemStorage, metricName, metricValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if metricType == model.MetricTypeGauge {
		err := gauge.StoreGauge(s.gaugeMemStorage, metricName, metricValue)
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
