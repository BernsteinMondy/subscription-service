package controller

import (
	"errors"
	srvc "github.com/BernsteinMondy/subscription-service/internal/service"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, srvc.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
