package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	if metrics.Value == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := fmt.Sprintf("%v", *metrics.Value)
	err = multistorage.SaveMetric(s.storage, metrics.ID, metrics.MType, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
