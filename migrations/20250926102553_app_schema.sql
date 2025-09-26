-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA app;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA app;
-- +goose StatementEnd
