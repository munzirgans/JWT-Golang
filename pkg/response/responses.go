package response

import (
	"encoding/json"
	"jwt/pkg/models"
	"net/http"
)

//JSON Response
func JSON(w http.ResponseWriter, statCode int, message string) {
	var response models.Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statCode)
	response.Message = message
	json.NewEncoder(w).Encode(response)
}

//Token Response
func Token(w http.ResponseWriter, statCode int, token string, message string) {
	var response models.RespToken
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statCode)
	response.Message = message
	response.Token = token
	json.NewEncoder(w).Encode(response)
}
