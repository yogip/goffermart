-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdrawals  DROP CONSTRAINT withdrawals_order_id_fkey;
ALTER TABLE withdrawals  DROP CONSTRAINT withdrawals_order_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdrawals ADD CONSTRAINT withdrawals_order_id_key FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE withdrawals ADD FOREIGN KEY (order_id) REFERENCES orders (id);
-- +goose StatementEnd

