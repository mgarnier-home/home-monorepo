package routes

import (
	"encoding/json"
	"mgarnier11/go/logger"
	"net/http"
)

type contextKey string

type Reponse struct {
	Message string `json:"message"`
}

func serializeAndSendResponse(w http.ResponseWriter, response interface{}, code int) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Error serializing response: %v", err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func sendOKResponse(w http.ResponseWriter, message string) {
	response := Reponse{Message: message}
	serializeAndSendResponse(w, response, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, message string, code int) {
	response := Reponse{Message: message}
	serializeAndSendResponse(w, response, code)
}
