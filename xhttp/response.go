package xhttp

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return NewError("error while processing response", http.StatusInternalServerError)
	}

	return nil
}
