package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/frieeze/autovoter/internal/signer"
	"github.com/frieeze/autovoter/internal/snapshot"
)

type App struct {
	config *config
	logger *slog.Logger
	signer *signer.Signer
	ss     *snapshot.Client

	ticker       *time.Ticker
	lastProposal string
	lastVote     time.Time
}

func New(logger *slog.Logger) (*App, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.Info("Config loaded", slog.Any("space", config.proposal.space), slog.Any("title", config.proposal.title), slog.Any("choice", config.proposal.choice))

	signer, err := signer.New(config.voter.address, config.voter.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	client := &snapshot.Client{
		Logger:    logger,
		HUB:       "https://hub.snapshot.org/graphql",
		Sequencer: "https://seq.snapshot.org/",
	}

	return &App{
		config: config,
		logger: logger,
		signer: signer,
		ss:     client,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.run(ctx)
	return nil
	/*
		a.ticker = time.NewTicker(time.Minute)
		defer a.ticker.Stop()
		a.run(ctx)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-a.ticker.C:
				a.run(ctx)
			}
		}
	*/
}

func (a *App) run(ctx context.Context) {
	/*
		if time.Since(a.lastVote) < 10*24*time.Hour ||
			time.Now().Weekday() != time.Thursday {
			return
		}
	*/

	proposal, choice, err := a.ss.GetProposal(ctx, a.config.proposal.space, a.config.proposal.title, a.config.proposal.choice)
	if err != nil {
		a.logger.Error("Failed to get proposal", slog.Any("error", err))
		return
	}

	if proposal == a.lastProposal {
		return
	}
	hasVoted, err := a.ss.HaveAlreadyVote(ctx, a.signer.Address, proposal)
	if err != nil {
		a.logger.Error("Failed to check if already voted", slog.Any("error", err))
		return
	}

	if hasVoted {
		a.logger.Info("Already voted", slog.Any("proposal", proposal))
		a.lastVote = time.Now()
		return
	}

	vote, sig, err := a.signer.Vote(choice, proposal, a.config.proposal.space)
	if err != nil {
		a.logger.Error("Failed to sign vote", slog.Any("error", err))
		return
	}
	fmt.Println(sig)

	if err := a.ss.SendVote(ctx, vote, sig); err != nil {
		a.logger.Error("Failed to send vote", slog.Any("error", err))
		return
	}

	a.logger.Info("Vote sent", slog.Any("proposal", proposal), slog.Any("choice", choice))

	a.lastProposal = proposal
	a.lastVote = time.Now()
}
