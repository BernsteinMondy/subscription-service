-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_subscriptions_all_filters
    ON app.subscriptions (service_name, user_id, start_date, end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX app.idx_subscriptions_all_filters;
-- +goose StatementEnd
