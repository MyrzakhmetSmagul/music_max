package musicmax

import (
	"encoding/json"
	"net/http"
)

type Description struct {
	Description string `json:"description"`
}

func DefaultResponse(w http.ResponseWriter, statusCode int) {
	errResp := Description{
		Description: http.StatusText(statusCode),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errResp)
}
