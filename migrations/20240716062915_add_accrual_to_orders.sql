-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN accrual INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN IF EXISTS accrual;
-- +goose StatementEnd
