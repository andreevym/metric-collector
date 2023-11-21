package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/multistorage"
)

// PostValueHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) PostValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ValueMetricContentType)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metrics := Metrics{}
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valStr, err := multistorage.GetMetric(s.storage, metrics.MType, metrics.ID)
	if err != nil || valStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resMetrics := Metrics{
		ID:    metrics.ID,
		MType: metrics.MType,
	}
	if resMetrics.MType == multistorage.MetricTypeGauge {
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resMetrics.Value = &v
	} else if resMetrics.MType == multistorage.MetricTypeCounter {
		v, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resMetrics.Delta = &v
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resBytes, err := json.Marshal(resMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(resBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
