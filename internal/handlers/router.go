package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	PathGetPing     = "/ping"
	PathPostUpdate  = "/update"
	PathPostUpdates = "/updates/"
	PathValue       = "/value"
	PathGetRoot     = "/"
)

func NewRouter(s *ServiceHandlers, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
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

	// Serve Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	r.Mount("/debug", chimiddleware.Profiler())

	return r
}
