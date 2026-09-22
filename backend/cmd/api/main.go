package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"developers-cv/backend/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.Initialize()
	if err != nil {
		panic(err)
	}

	if err := application.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		application.Logger().Error("application stopped with error", zap.Error(err))
		os.Exit(1)
	}
}
