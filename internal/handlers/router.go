package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(s *ServiceHandlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.PostUpdateHandler)
	r.Get("/value/{metricType}/{metricName}", s.GetValueHandler)
	r.Post("/update", s.PostUpdateHandler)
	r.Post("/update/", s.PostUpdateHandler)
	r.Post("/value", s.PostValueHandler)
	r.Post("/value/", s.PostValueHandler)
	return r
}
