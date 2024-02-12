package application

import (
	"context"
	"fmt"
	"github.com/sushil-cmd-r/order-api/db"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	db     db.Database
	config Config
}

func NewApp(config Config, database db.Database) *App {

	app := &App{
		db:     database,
		config: config,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	err := a.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer a.db.Close()

	a.loadRoutes()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	fmt.Printf("Starting server on port %d...\n", a.config.ServerPort)
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
