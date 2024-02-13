package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/sushil-cmd-r/order-api/db"
	"gorm.io/gorm"
	"os"
)

type PostgresRepo struct {
	Client *gorm.DB
}

func NewPostgresRepo(pdb db.Database) *PostgresRepo {
	client, ok := pdb.(*db.PostgresDb)
	if !ok {
		// should be unreachable
		fmt.Println("unable to cast database")
		os.Exit(1)
	}

	err := client.GetDB().AutoMigrate(&Order{}, &LineItem{})
	if err != nil {
		fmt.Println("failed to run migrations: %w", err)
		os.Exit(1)
	}

	return &PostgresRepo{
		Client: client.GetDB(),
	}
}

func (p *PostgresRepo) Insert(_ context.Context, order Order) error {
	res := p.Client.Create(&order)
	if res.Error != nil {
		return fmt.Errorf("unable to insert order: %w", res.Error)
	}
	return nil
}

func (p *PostgresRepo) FindByID(_ context.Context, id string) (Order, error) {
	var order Order

	res := p.Client.Preload("LineItems").Where("order_id = ?", id).First(&order)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Order{}, ErrNotExist
	}
	if res.Error != nil {
		return Order{}, fmt.Errorf("unable to find by id: %w", res.Error)
	}
	return order, nil
}

func (p *PostgresRepo) DeleteById(_ context.Context, id string) error {
	return nil
}

func (p *PostgresRepo) Update(_ context.Context, order Order) error {
	p.Client.Save(&order)
	return nil
}

func (p *PostgresRepo) FindAll(_ context.Context, page FindAllPage) (FindResult, error) {

	limit := int(page.Size)
	orders := make([]Order, 0)
	res := p.Client.Preload("LineItems").Offset(int(page.Offset)).Limit(limit).Find(&orders)

	if res.Error != nil {
		return FindResult{}, fmt.Errorf("unable to get orders: %w", res.Error)
	}

	var nxtCursor uint64
	if len(orders) == limit {
		nxtCursor = page.Offset + uint64(len(orders))
	}

	return FindResult{
		Orders: orders,
		Cursor: nxtCursor,
	}, nil
}
