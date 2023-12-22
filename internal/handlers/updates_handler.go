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
