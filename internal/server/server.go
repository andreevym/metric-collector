package server

import (
	"log"
	"net/http"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/go-chi/chi"
)

func StartServer() {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	s := handlers.NewServer(counterMemStorage, gaugeMemStorage)

	r := chi.NewRouter()
	r.Handle(
		"/update/{metricType}/{metricName}/{metricValue}",
		s.UpdateMetricHandler(),
	)
	r.Get(
		"/value/{metricType}/{metricName}",
		s.GetMetricByTypeAndNameHandler(),
	)
	log.Fatal(http.ListenAndServe(":8080", r))
}
