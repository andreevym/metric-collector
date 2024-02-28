package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"go.uber.org/zap"
)

const (
	UpdatesMetricContentType = "application/json"
)

// PostUpdatesHandler method for bulk insert or update of metrics.
// This endpoint is used to bulk insert or update metric values by sending a POST request with a JSON array of metrics.
// @Summary Bulk insert or update metrics
// @Description Bulk inserts or updates metric values.
// This endpoint accepts a POST request with a JSON array of metrics.
// Each metric should have an ID, type, and either delta (for counter type) or value (for gauge type).
// Supported metric types are 'gauge' and 'counter'.
// @Accept json
// @Produce json
// @Param metrics body []storage.Metric true "Array of metrics to insert or update"
// @Success 200 {string} string "Metrics inserted or updated successfully"
// @Failure 400 {string} string "Bad request. Invalid JSON payload or metric parameters"
// @Router /updates [post]
func (s ServiceHandlers) PostUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdatesMetricContentType)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Logger().Error("error", zap.Error(err))
		return
	}

	var metrics []*storage.Metric
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		logger.Logger().Error("err", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage.SaveAllMetric(r.Context(), s.storage, metrics)
}
