package entity

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ServiceName string
	Price       int32
	StartDate   time.Time
	EndDate     time.Time
}

type CreateSubscriptionData struct {
	UserID      uuid.UUID
	ServiceName string
	Price       int32
	StartDate   time.Time
	EndDate     time.Time
}

type UpdateSubscriptionData struct {
	ServiceName string
	Price       int32
	StartDate   time.Time
	EndDate     time.Time
}

type GetSubscriptionsFilter struct {
	UserID      uuid.UUID
	ServiceName string
	StartDate   time.Time
	EndDate     time.Time
}
