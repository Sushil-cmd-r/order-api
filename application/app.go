package application

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func NewApp() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to order: %w", err)
	}

	defer func() {
		if err = a.rdb.Close(); err != nil {
			fmt.Println("failed to close order", err)
		}
	}()

	fmt.Println("Starting server...")

	ch := make(chan error, 1)

	go func(ch chan error) {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}(ch)

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Shutting down gracefully...")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
