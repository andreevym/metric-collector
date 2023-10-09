package handlers

import "github.com/go-chi/chi"

func Router(s Server) chi.Router {
	r := chi.NewRouter()
	r.Handle(
		"/update/{metricType}/{metricName}/{metricValue}",
		s.UpdateMetricHandler(),
	)
	r.Get(
		"/value/{metricType}/{metricName}",
		s.GetMetricByTypeAndNameHandler(),
	)
	return r
}
