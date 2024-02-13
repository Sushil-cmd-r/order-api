package main

import (
	"context"
	"fmt"
	"github.com/sushil-cmd-r/order-api/application"
	"github.com/sushil-cmd-r/order-api/db"
	"os"
	"os/signal"
)

func main() {
	config := application.LoadConfig()
	database := db.NewPostgresDb(config.PostgresAddr)

	app := application.NewApp(config, database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
