package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/multistorage"
)

// ValueMetricByTypeAndNameHandler method return metric value by metric type and metric name
// example request url: http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s ServiceHandlers) ValueMetricByTypeAndNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ValueMetricContentType)

	bytes, _ := io.ReadAll(r.Body)
	metrics := Metrics{}
	err := json.Unmarshal(bytes, &metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valStr, err := multistorage.GetMetric(s.storage, metrics.MType, metrics.ID)
	if err != nil || valStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics := Metrics{
			ID:    metrics.ID,
			MType: metrics.MType,
			Delta: nil,
			Value: &v,
		}
		bytes, err := json.Marshal(metrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = io.WriteString(w, string(bytes))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
