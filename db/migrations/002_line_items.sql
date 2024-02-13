-- +goose Up
CREATE TABLE line_items (
    item_id uuid primary key ,
    quantity INT,
    price INT,
    order_id varchar,

    constraint orders_line_items foreign key (order_id)
    references orders(order_id)  on delete cascade
);

-- +goose Down
DROP TABLE line_items;