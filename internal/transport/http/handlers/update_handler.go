package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	UpdateMetricContentType = "application/json"
	ValueMetricContentType  = "application/json"
)

// PostUpdateHandler method for insert or update metrics.
// This endpoint is used to insert or update metric values by sending a POST request with the metric ID, type, and value.
// @Summary Insert or update metric value
// @Description Inserts or updates the value of a metric specified by its type, name, and value.
// This endpoint accepts a POST request with the metric ID, type, and value as path parameters.
// Supported metric types are 'gauge' and 'counter'.
// @Param metricType path string true "Type of the metric ('gauge' or 'counter')"
// @Param metricName path string true "Name of the metric"
// @Param metricValue path number true "Value of the metric"
// @Produce json
// @Success 200 {object} store.Metric "Metric value inserted or updated successfully"
// @Failure 400 {string} string "Bad request. Invalid metric parameters or JSON payload"
// @Router /update/{metricType}/{metricName}/{metricValue} [post]
func (s ServiceHandlers) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", UpdateMetricContentType)

	metric, err := buildMetricByRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respMetric, err := s.controller.Update(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(&respMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func buildMetricByRequest(r *http.Request) (*store.Metric, error) {
	metric, err := buildMetricByBody(r.Body)
	if err == nil && metric != nil {
		return metric, nil
	}

	return BuildMetricByChiParam(r)
}

func buildMetricByBody(body io.ReadCloser) (*store.Metric, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		logger.Logger().Error("error", zap.Error(err))
	}

	if len(bytes) == 0 {
		return nil, errors.New("body len is empty")
	}

	metric := &store.Metric{}
	err = json.Unmarshal(bytes, &metric)
	return metric, err
}

func BuildMetricByChiParam(r *http.Request) (*store.Metric, error) {
	metric := &store.Metric{}
	metric.MType = chi.URLParam(r, "metricType")
	metric.ID = chi.URLParam(r, "metricName")
	v := chi.URLParam(r, "metricValue")
	if metric.MType == store.MTypeCounter {
		delta, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		metric.Delta = &delta
	} else if metric.MType == store.MTypeGauge {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		metric.Value = &value
	} else {
		return nil, errors.New("unknown type")
	}

	return metric, nil
}

const (
	countRequestURIParam    = 5
	indexMetricName         = 3
	indexMetricType         = 2
	indexMetricValueOrDelta = 4
)

func BuildMetricBySplitParam(r *http.Request) (*store.Metric, error) {
	split := make([]string, countRequestURIParam)
	copy(split, strings.Split(r.URL.Path, "/"))

	if len(split) < countRequestURIParam {
		return nil, errors.New("unknown size param")
	}

	m := &store.Metric{
		MType: split[indexMetricType],
		ID:    split[indexMetricName],
	}

	valueOrDelta := split[indexMetricValueOrDelta]

	switch m.MType {
	case store.MTypeCounter:
		delta, err := strconv.ParseInt(valueOrDelta, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse delta for metric: %w", err)
		}
		m.Delta = &delta
	case store.MTypeGauge:
		value, err := strconv.ParseFloat(valueOrDelta, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value for metric: %w", err)
		}
		m.Value = &value
	default:
		return nil, errors.New("unknown type")
	}

	return m, nil
}
