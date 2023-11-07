package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/go-chi/chi/v5"
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

// PostUpdateHandler method for insert or update metrics
// example request url: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s ServiceHandlers) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdateMetricContentType)

	bytes, _ := io.ReadAll(r.Body)

	var metricName string
	var metricType string
	var metricValue string

	if len(bytes) > 0 {
		metrics := Metrics{}
		err := json.Unmarshal(bytes, &metrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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

		metricType = metrics.MType
		metricName = metrics.ID
	} else {
		metricType = chi.URLParam(r, "metricType")
		metricName = chi.URLParam(r, "metricName")
		metricValue = chi.URLParam(r, "metricValue")
	}

	newVal, err := multistorage.SaveMetric(s.storage, metricName, metricType, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := Metrics{
		ID:    metricName,
		MType: metricType,
	}
	if resp.MType == multistorage.MetricTypeGauge {
		v, err := strconv.ParseFloat(newVal, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Value = &v
	} else if resp.MType == multistorage.MetricTypeCounter {
		v, err := strconv.ParseInt(newVal, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Delta = &v
	}
	bytesResp, err := json.Marshal(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(bytesResp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
