package db

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDb struct {
	dataSrc string
	client  *gorm.DB
}

func NewPostgresDb(dataSrc string) *PostgresDb {
	return &PostgresDb{
		dataSrc: dataSrc,
	}
}

func (p *PostgresDb) Connect(context.Context) error {
	db, err := gorm.Open(postgres.Open(p.dataSrc), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	p.client = db
	return nil
}

func (p *PostgresDb) GetDB() *gorm.DB {
	return p.client
}

func (p *PostgresDb) Close() {
	db, _ := p.client.DB()
	_ = db.Close()
}
