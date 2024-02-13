package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sushil-cmd-r/order-api/db"
	"os"
	"time"
)

type PostgresRepo struct {
	Client *pgxpool.Pool
}

func NewPostgresRepo(pgdb db.Database) *PostgresRepo {
	client, ok := pgdb.(*db.PostgresDb)
	if !ok {
		// should be unreachable
		fmt.Println("unable to cast database")
		os.Exit(1)
	}

	return &PostgresRepo{
		Client: client.GetDB(),
	}
}

func (p *PostgresRepo) Insert(ctx context.Context, order Order) error {
	tx, err := p.Client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q1 := `INSERT INTO orders (order_id, customer_id, created_at) VALUES ($1, $2, $3);`
	ct, err := tx.Exec(ctx, q1, order.OrderID, order.CustomerID, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to insert order: %w", err)
	}
	if ct.RowsAffected() != 1 {
		return fmt.Errorf("invalid insert")
	}

	lineItems := order.LineItems
	cc, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"line_items"},
		[]string{"item_id", "quantity", "price", "order_id"},
		pgx.CopyFromSlice(len(lineItems), insertLineItems(order.OrderID, lineItems)),
	)
	if err != nil {
		return fmt.Errorf("unable insert line items: %w", err)
	}
	if int(cc) != len(lineItems) {
		return fmt.Errorf("invalid insert")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit insert: %w", err)
	}

	return nil
}

func insertLineItems(orderId string, lineItems []LineItem) func(int) ([]any, error) {
	return func(i int) ([]any, error) {
		return []any{
			lineItems[i].ItemID,
			lineItems[i].Quantity,
			lineItems[i].Price,
			orderId,
		}, nil
	}
}

func (p *PostgresRepo) FindByID(ctx context.Context, id string) (Order, error) {
	tx, err := p.Client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Order{}, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var order Order
	q1 := `SELECT * from orders where order_id = $1;`
	err = tx.QueryRow(ctx, q1, id).Scan(&order)

	if errors.Is(pgx.ErrNoRows, err) {
		return order, ErrNotExist
	} else if err != nil {
		return order, fmt.Errorf("unable to get order by id: %w", err)
	}

	q2 := `SELECT * from line_items where order_id = $1;`
	rows, err := tx.Query(ctx, q2, id)
	if err != nil {
		return order, fmt.Errorf("unable to get line items by id: %w", err)
	}

	lineItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[LineItem])
	if err != nil {
		return order, fmt.Errorf("unable to cast to line item: %w", err)
	}

	order.LineItems = lineItems

	return order, nil
}

func (p *PostgresRepo) DeleteById(ctx context.Context, id string) error {
	tx, err := p.Client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q1 := `DELETE FROM orders where order_id = $1`
	ct, err := tx.Exec(ctx, q1, id)

	if err != nil {
		return fmt.Errorf("unable to delete order by id: %w", err)
	}
	if ct.RowsAffected() != 1 {
		return ErrNotExist
	}

	return nil
}

func (p *PostgresRepo) Update(ctx context.Context, order Order) error {
	return nil
}

func (p *PostgresRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {

	tx, err := p.Client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return FindResult{}, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ts := time.Unix(int64(page.Offset), 0).UTC()
	q1 := `SELECT * from orders where created_at > $1 limit $2;`
	rows, err := tx.Query(ctx, q1, ts, page.Size)
	if err != nil {
		return FindResult{}, fmt.Errorf("unable to query orders: %w", err)
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[Order])
	if err != nil {
		return FindResult{}, fmt.Errorf("unable to collect orders: %w", err)
	}

	if len(orders) == -1 {
		return FindResult{
			Orders: nil,
		}, nil
	}
	q2 := `SELECT * from line_items where order_id = $1;`
	for _, o := range orders {
		rows, err := tx.Query(ctx, q2, o.OrderID)
		if err != nil {
			return FindResult{}, fmt.Errorf("unable to get line items by id: %w", err)
		}

		lineItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[LineItem])
		if err != nil {
			return FindResult{}, fmt.Errorf("unable to cast to line item: %w", err)
		}

		o.LineItems = lineItems
	}

	res := FindResult{
		Orders: orders,
		Cursor: uint64(orders[len(orders)-1].CreatedAt.Unix()),
	}

	return res, nil
}
