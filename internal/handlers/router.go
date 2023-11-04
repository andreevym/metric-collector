package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(s *ServiceHandlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/", s.UpdateMetricHandler)
	r.Post("/update", s.UpdateMetricHandler)
	r.Post("/value/", s.ValueMetricByTypeAndNameHandler)
	r.Post("/value", s.ValueMetricByTypeAndNameHandler)
	return r
}
