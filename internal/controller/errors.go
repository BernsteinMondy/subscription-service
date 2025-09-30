package controller

import (
	"errors"
	srvc "github.com/BernsteinMondy/subscription-service/internal/service"
	"log/slog"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, srvc.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("resource not found"))
		return
	default:
		slog.Error("unexpected internal error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
