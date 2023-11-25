package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/storage"
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

	m := storage.Metric{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundMetric, err := s.storage.Read(r.Context(), m.ID)
	if err != nil || foundMetric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err = json.Marshal(foundMetric)
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
