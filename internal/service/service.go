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
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	GetAllSubscriptionsFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) ([]entity.Subscription, error)

	CreateSubscription(ctx context.Context, subscription *entity.Subscription) (uuid.UUID, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, data *entity.UpdateSubscriptionData) error
	DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	sub, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get subscription by id: %w", err)
	}

	return sub, nil
}

func (s *service) NewSubscription(ctx context.Context, data *entity.CreateSubscriptionData) (uuid.UUID, error) {
	sub := &entity.Subscription{
		ID:          uuid.New(),
		UserID:      data.UserID,
		ServiceName: data.ServiceName,
		Price:       data.Price,
		StartDate:   data.StartDate,
		EndDate:     data.EndDate,
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

func (s *service) UpdateSubscription(ctx context.Context, id uuid.UUID, data *entity.UpdateSubscriptionData) error {
	err := s.repo.UpdateSubscription(ctx, id, data)
	if err != nil {
		if errors.Is(err, repo.ErrRepoNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("repo: update subscription: %w", err)
	}

	return nil
}

func (s *service) GetSubscriptionsTotalSumFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) (int32, error) {
	subs, err := s.repo.GetAllSubscriptionsFilter(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("repo: get all subscriptions with filter: %w", err)
	}

	totalPrice := int32(0)
	for _, sub := range subs {
		totalPrice += sub.Price
	}

	return totalPrice, nil
}

func (s *service) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	subs, err := s.repo.GetAllSubscriptionsFilter(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("repo: get all subscriptions: %w", err)
	}

	return subs, nil
}
