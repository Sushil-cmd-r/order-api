package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDb struct {
	dataSrc string
	client  *pgxpool.Pool
}

func NewPostgresDb(dataSrc string) *PostgresDb {
	return &PostgresDb{
		dataSrc: dataSrc,
	}
}

func (p *PostgresDb) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, p.dataSrc)
	if err != nil {
		return err
	}

	if pool == nil {
		return errors.New("unable to create database")
	}
	p.client = pool
	return nil
}

func (p *PostgresDb) GetDB() *pgxpool.Pool {
	return p.client
}

func (p *PostgresDb) Close() {
	p.client.Close()
}
