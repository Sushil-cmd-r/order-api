package db

import "context"

type Database interface {
	Connect(ctx context.Context) error
	Close()
}
