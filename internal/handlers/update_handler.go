package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/multistorage"
)

const (
	UpdateMetricContentType = "application/json"
	ValueMetricContentType  = "application/json"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

// UpdateMetricHandler method for insert or update metrics
// example request url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s ServiceHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdateMetricContentType)

	bytes, _ := io.ReadAll(r.Body)
	metrics := Metrics{}
	err := json.Unmarshal(bytes, &metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metricValue string
	if metrics.MType == multistorage.MetricTypeGauge {
		if metrics.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricValue = fmt.Sprintf("%v", *metrics.Value)
	} else if metrics.MType == multistorage.MetricTypeCounter {
		if metrics.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricValue = strconv.FormatInt(*metrics.Delta, 10)
	}

	newVal, err := multistorage.SaveMetric(s.storage, metrics.ID, metrics.MType, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metrics.MType == multistorage.MetricTypeGauge {
		v, err := strconv.ParseFloat(newVal, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics.Value = &v
	} else if metrics.MType == multistorage.MetricTypeCounter {
		v, err := strconv.ParseInt(newVal, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics.Delta = &v
	}

	w.WriteHeader(http.StatusOK)
}
