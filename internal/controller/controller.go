package controller

import (
	"context"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type service interface {
	CancelSubscription(ctx context.Context, id uuid.UUID) error
	NewSubscription(ctx context.Context, data *entity.CreateSubscriptionData) (uuid.UUID, error)
	GetSubscriptionsTotalSumFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) ([]entity.Subscription, error)
}

type controller struct {
	service service
}

func New(srvc service) *controller {
	return &controller{
		service: srvc,
	}
}
