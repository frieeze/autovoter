package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/frieeze/autovoter/internal/signer"
	"github.com/frieeze/autovoter/internal/snapshot"
	"github.com/joho/godotenv"
)

type App struct {
	logger *slog.Logger
	signer *signer.Signer
	client *snapshot.Client
}

func New(logger *slog.Logger) (*App, error) {
	godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	signer, err := signer.New(config.Voter.Address, config.Voter.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	client := &snapshot.Client{
		HUB:       "https://hub.snapshot.org",
		Sequencer: "https://seq.snapshot.org",
	}

	return &App{
		logger: logger,
		signer: signer,
		client: client,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	return nil
}
