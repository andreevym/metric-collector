package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

// PostValueHandler method returns metric value by metric type and metric name.
// This endpoint is used to retrieve metric values by sending a POST request with JSON payload.
// @Summary Retrieve metric value by type and name
// @Description Retrieves the value of a metric specified by its type and name.
// This endpoint accepts a JSON payload containing the metric ID and type.
// Supported metric types are 'gauge' and 'counter'.
// @Param metricType path string true "Type of the metric ('gauge' or 'counter')"
// @Param metricName path string true "Name of the metric"
// @Accept json
// @Produce json
// @Success 200 {object} store.Metric "Metric value retrieved successfully"
// @Failure 400 {string} string "Bad request. Invalid JSON payload"
// @Failure 404 {string} string "Metric value not found"
// @Router /value/{metricType}/{metricName} [post]
func (s ServiceHandlers) PostValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ValueMetricContentType)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := store.Metric{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundMetric, err := s.storage.Read(r.Context(), m.ID, m.MType)
	if err != nil || foundMetric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err = json.Marshal(foundMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = io.WriteString(w, string(bytes))
	if err != nil {
		logger.Logger().Error("value can't be written", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
