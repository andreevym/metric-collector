package handlers

import "net/http"

// GetPingHandler ping database
// @Summary Ping database
// @Description Pings the database to check its connectivity
// @Success 200 {string} string "Database ping successful"
// @Failure 500 {string} string "Internal server error"
// @Router /ping [get]
func (s ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	if s.dbClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := s.dbClient.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
