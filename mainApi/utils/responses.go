package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type DefaultResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	ToJson(w, struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})
}

func ToJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logrus.WithError(err).Warn("Error writing response")
	}
}
