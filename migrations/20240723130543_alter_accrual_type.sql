-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE NUMERIC(8, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE int;
-- +goose StatementEnd
