package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(s *ServiceHandlers, middlewares ...func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.PostUpdateHandler)
	r.Get("/value/{metricType}/{metricName}", s.GetValueHandler)
	r.Get("/ping", s.GetPingHandler)
	r.Post("/update", s.PostUpdateHandler)
	r.Post("/updates/", s.PostUpdatesHandler)
	r.Post("/update/", s.PostUpdateHandler)
	r.Post("/value", s.PostValueHandler)
	r.Post("/value/", s.PostValueHandler)
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	})
	return r
}
