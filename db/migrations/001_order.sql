-- +goose Up
CREATE TABLE orders (
    order_id varchar primary key ,
    customer_id uuid,
    created_at timestamp,
    shipped_at timestamp,
    completed_at timestamp
);

-- +goose Down
DROP TABLE orders;