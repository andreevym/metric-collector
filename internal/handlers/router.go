package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(s *ServiceHandlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.UpdateMetricHandler)
	r.Get("/value/{metricType}/{metricName}", s.GetValueHandler)
	r.Post("/update", s.UpdateMetricHandler)
	r.Post("/update/", s.UpdateMetricHandler)
	r.Post("/value", s.ValueMetricByTypeAndNameHandler)
	r.Post("/value/", s.ValueMetricByTypeAndNameHandler)
	return r
}
