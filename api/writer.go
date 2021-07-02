package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorResponse represents a error response
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// WriteError handles writing HTTP error responses
func WriteError(w http.ResponseWriter, statusCode int, log *logrus.Entry, errs ...error) {
	er := ErrorResponse{Errors: []string{}}

	for _, err := range errs {
		er.Errors = append(er.Errors, err.Error())
	}

	resp, err := json.Marshal(er)
	if err != nil {
		log.WithError(err).Error("failed to marshal error response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, werr := w.Write(resp); werr != nil {
		log.WithError(err).Error("failed to write response")
	}
}

// WriteJSON handler writing HTTP JSON responses
func WriteJSON(w http.ResponseWriter, log *logrus.Entry, out interface{}) {
	resp, err := json.Marshal(out)
	if err != nil {
		log.WithError(err).Error("failed to marshal error JSON response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, werr := w.Write(resp); werr != nil {
		log.WithError(err).Error("failed to write JSON response")
	}
}
