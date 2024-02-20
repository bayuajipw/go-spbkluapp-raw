package maincontroller

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, message string) {
	RespondWithJSON(w, map[string]interface{}{"success": false, "message": message})
}

func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
