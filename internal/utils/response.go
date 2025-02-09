package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // Optional field for data
}

func SendJSONResponse(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JsonResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
