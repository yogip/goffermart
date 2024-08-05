-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdrawals ALTER COLUMN sum TYPE NUMERIC(8, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdrawals ALTER COLUMN sum TYPE bigint;
-- +goose StatementEnd
