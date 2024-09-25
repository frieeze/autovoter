package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/frieeze/autovoter/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app, err := app.New(logger)
	if err != nil {
		logger.Error("Failed to create app", slog.Any("error", err))
		return
	}

	if err := app.Start(ctx); err != nil {
		logger.Error("Failed to start app", slog.Any("error", err))
	}

	logger.Info("App stopped gracefully")
}
