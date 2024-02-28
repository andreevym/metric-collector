package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

// GetValueHandler method return metric value by metric type and metric name
// @Summary Retrieve metric value by type and name
// @Description Retrieves the value of a metric specified by its type and name.
// Supported metric types are 'gauge' and 'counter'.
// @Param metricType path string true "Type of the metric ('gauge' or 'counter')"
// @Param metricName path string true "Name of the metric"
// @Success 200 {string} string "Metric value retrieved successfully"
// @Failure 400 {string} string "Bad request. Either metric type is unsupported or value is missing"
// @Failure 404 {string} string "Metric value not found"
// @Router /value/{metricType}/{metricName} [get]
func (s ServiceHandlers) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	if metricType != storage.MTypeGauge && metricType != storage.MTypeCounter {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricName := chi.URLParam(r, "metricName")
	v, err := s.storage.Read(r.Context(), metricName, metricType)
	if err != nil || v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch v.MType {
	case storage.MTypeCounter:
		if v.Delta == nil {
			logger.Logger().Error("delta can't be nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strconv.FormatInt(*v.Delta, 10)
		_, err = io.WriteString(w, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case storage.MTypeGauge:
		if v.Value == nil {
			logger.Logger().Error("value can't be nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strconv.FormatFloat(*v.Value, 'f', -1, 64)
		_, err = io.WriteString(w, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
