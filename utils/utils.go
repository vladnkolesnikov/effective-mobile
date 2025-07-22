package utils

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]any

func WriteResponse(w http.ResponseWriter, status int, envelope Envelope) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(envelope); err != nil {
		return err
	}

	return nil
}
