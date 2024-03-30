package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	PathGetPing     = "/ping"
	PathPostUpdate  = "/update"
	PathPostUpdates = "/updates/"
	PathValue       = "/value"
	PathGetRoot     = "/"
)

func NewRouter(
	s *ServiceHandlers,
	middlewares ...func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares...)

	r.Get(PathGetPing, s.GetPingHandler)

	r.Post(PathPostUpdates, s.PostUpdatesHandler)

	r.Post(PathPostUpdate, s.PostUpdateHandler)
	r.Post(PathPostUpdate+"/", s.PostUpdateHandler)
	r.Post(PathPostUpdate+"/{metricType}/{metricName}/{metricValue}", s.PostUpdateHandler)

	r.Post(PathValue, s.PostValueHandler)
	r.Post(PathValue+"/", s.PostValueHandler)
	r.Get(PathValue+"/{metricType}/{metricName}", s.GetValueHandler)

	r.Get(PathGetRoot, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	})
	return r
}
