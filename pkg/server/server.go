package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/subscription-service/pkg/server/http"
)

type Option func(server Server)

type Type uint8

const (
	HTTP Type = iota + 1
	GRPC
)

type Server interface {
	Launch() error
	Shutdown(ctx context.Context) error
}

func New(serverType Type, opts ...Option) (Server, error) {
	var s Server

	switch serverType {
	case HTTP:
		s = http.NewServer()
	case GRPC:
		// TODO: Maybe add gRPC server
		return nil, errors.New("not implemented")
	default:
		return nil, fmt.Errorf("unknown server type: %v", serverType)
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
