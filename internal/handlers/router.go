package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(s *ServiceHandlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.UpdateMetricHandler)
	r.Get("/value/{metricType}/{metricName}", s.GetMetricByTypeAndNameHandler)
	return r
}
