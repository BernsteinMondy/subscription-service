package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	repo "github.com/BernsteinMondy/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type repository interface {
	CreateSubscription(ctx context.Context, subscription *entity.Subscription) (uuid.UUID, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error
	GetAllSubscriptionsFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) ([]entity.Subscription, error)
}

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) NewSubscription(ctx context.Context, data *entity.CreateSubscriptionData) (uuid.UUID, error) {
	sub := &entity.Subscription{
		ID:          uuid.New(),
		UserID:      data.UserID,
		ServiceName: data.ServiceName,
		Price:       data.Price,
		StartDate:   data.StartDate,
	}

	id, err := s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repo: create subscription: %w", err)
	}

	return id, nil
}

func (s *service) CancelSubscription(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrRepoNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("repo: delete subscription: %w", err)
	}

	return nil
}

func (s *service) GetSubscriptionsTotalSumFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) ([]entity.Subscription, error) {
	subs, err := s.repo.GetAllSubscriptionsFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("repo: get all subscriptions with filter: %w", err)
	}

	return subs, nil
}
