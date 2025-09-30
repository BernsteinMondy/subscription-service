package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
	"strings"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateSubscription(ctx context.Context, subscription *entity.Subscription) (uuid.UUID, error) {
	const query = `INSERT INTO app.subscriptions (id, user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, subscription.ID, subscription.UserID, subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate)
	if err != nil {
		return uuid.Nil, err
	}

	return subscription.ID, nil
}

func (r *repository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	const query = `SELECT user_id, service_name, price, start_date, end_date FROM app.subscriptions WHERE id = $1`

	res := entity.Subscription{
		ID: id,
	}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&res.UserID, &res.ServiceName, &res.Price, &res.StartDate, &res.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &res, nil
}

func (r *repository) DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM app.subscriptions WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRepoNotFound
	}

	return nil
}

func (r *repository) UpdateSubscription(ctx context.Context, id uuid.UUID, data *entity.UpdateSubscriptionData) error {
	const query = `UPDATE app.subscriptions SET price = $1, service_name = $2, start_date = $3, end_date = $4 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, data.Price, data.EndDate, id)
	if err != nil {
		return fmt.Errorf("exec sql query: %w", err)
	}

	return nil
}

func (r *repository) GetAllSubscriptionsFilter(ctx context.Context, filter *entity.GetSubscriptionsFilter) (_ []entity.Subscription, err error) {
	var (
		queryBuilder strings.Builder
		args         []interface{}
		conditions   []string
	)

	queryBuilder.WriteString(`SELECT id, user_id, service_name, price, start_date, end_date FROM app.subscriptions`)

	// Собираем условия
	if filter.ServiceName != "" {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(args)+1))
		args = append(args, filter.ServiceName)
	}

	if filter.UserID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, filter.UserID)
	}

	if !filter.StartDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", len(args)+1))
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", len(args)+1))
		args = append(args, filter.EndDate)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	query := queryBuilder.String()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close rows: %w", closeErr))
		}
	}()

	var subscriptions []entity.Subscription

	for rows.Next() {
		var subscription entity.Subscription
		err = rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.StartDate,
			&subscription.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return subscriptions, nil
}
