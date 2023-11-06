package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/go-chi/chi/v5"
)

// ValueMetricByTypeAndNameHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) ValueMetricByTypeAndNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ValueMetricContentType)

	bytes, _ := io.ReadAll(r.Body)

	var metricType string
	var metricName string
	if len(bytes) > 0 {
		metrics := Metrics{}
		err := json.Unmarshal(bytes, &metrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricType = metrics.MType
		metricName = metrics.ID
	} else {
		metricType = chi.URLParam(r, "metricType")
		metricName = chi.URLParam(r, "metricName")
	}

	valStr, err := multistorage.GetMetric(s.storage, metricType, metricName)
	if err != nil || valStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resMetrics := Metrics{
		ID:    metricName,
		MType: metricType,
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

	bytes, err = json.Marshal(resMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
