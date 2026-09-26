package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"developers-cv/backend/internal/app"
	"developers-cv/backend/internal/platform/healthcheck"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := healthcheck.Check(context.Background()); err != nil {
			os.Exit(1)
		}
		return
	}

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
