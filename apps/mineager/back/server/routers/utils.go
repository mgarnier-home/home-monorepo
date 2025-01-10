package routers

import (
	"encoding/json"
	"mgarnier11/go/logger"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type RouterUtils struct {
	logger *logger.Logger
}

func (utils *RouterUtils) serializeAndSendResponse(w http.ResponseWriter, r *http.Request, response interface{}, code int) {
	jsonResponse, err := json.Marshal(response)

	if err != nil {
		utils.logger.Errorf("Error serializing response: %v", err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)

	utils.logger.Infof("%s %s Response %d", r.Method, r.URL.Path, code)
}

func (utils *RouterUtils) sendOKResponse(w http.ResponseWriter, r *http.Request, message string) {
	response := Response{Message: message}
	utils.serializeAndSendResponse(w, r, response, http.StatusOK)
}

func (utils *RouterUtils) sendErrorResponse(w http.ResponseWriter, r *http.Request, message string, code int) {
	response := Response{Message: message}
	utils.serializeAndSendResponse(w, r, response, code)
}
