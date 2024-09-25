package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.WithCancel(context.Background(), os.Interrupt))
	defer cancel()

	godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		logger.Error("Failed to load config", slog.Any("error", err))
	}

}
