package handlers

import "net/http"

// GetPingHandler ping database
func (s ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	err := s.dbClient.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
