package musicmax

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Description struct {
	Description string `json:"description"`
}

var ErrBadRequest = errors.New("bad request")
var ErrInternalServerError = errors.New("internal server error")

func DefaultResponse(w http.ResponseWriter, statusCode int) {
	errResp := Description{
		Description: http.StatusText(statusCode),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errResp)
}
