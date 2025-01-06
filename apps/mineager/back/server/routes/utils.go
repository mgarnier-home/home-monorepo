package routes

import (
	"encoding/json"
	"mgarnier11/go/logger"
	"net/http"
)

type contextKey string

func serializeAndSendResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Error serializing response: %v", err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
