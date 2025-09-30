package controller

import (
	"context"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type service interface {
	GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error)
	GetSubscriptionsTotalSumFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) (int32, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)

	CancelSubscription(ctx context.Context, id uuid.UUID) error
	NewSubscription(ctx context.Context, data *entity.CreateSubscriptionData) (uuid.UUID, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, data *entity.UpdateSubscriptionData) error
}

type controller struct {
	service service
}

func New(srvc service) *controller {
	return &controller{
		service: srvc,
	}
}
