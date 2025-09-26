-- +goose Up
-- +goose StatementBegin
CREATE TABLE app.subscriptions
(
    id           uuid        NOT NULL PRIMARY KEY,
    user_id      uuid        NOT NULL,
    service_name text        NOT NULL,
    price        integer     NOT NULL,
    start_date   timestamptz NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE app.subscriptions;
-- +goose StatementEnd
