package main

import (
	"context"
	"fmt"
	"github.com/sushil-cmd-r/order-api/application"
	"os"
	"os/signal"
)

func main() {
	app := application.NewApp(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
