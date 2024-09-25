package app

import (
	"errors"
	"os"
)

// Voter represents the delegate who will vote
// is fetched from the environment variables
type Voter struct {
	Address    string
	PrivateKey string
}

type Proposal struct {
	Choice string
	Space  string
	Title  string
}

type Config struct {
	Voter    Voter
	Proposal Proposal
}

var (
	errNoVoterAddress    = errors.New("VOTER_ADDRESS is not set")
	errNoVoterPrivateKey = errors.New("VOTER_PRIVATE_KEY is not set")
	errNoChoice          = errors.New("PROPOSAL_CHOICE is not set")
	errNoSpace           = errors.New("PROPOSAL_SPACE is not set")
	errNoProposalTitle   = errors.New("PROPOSAL_TITLE is not set")
)

func loadConfig() (*Config, error) {
	// Load the configuration from the environment variables
	address, exists := os.LookupEnv("VOTER_ADDRESS")
	if !exists {
		return nil, errNoVoterAddress
	}

	privateKey, exists := os.LookupEnv("VOTER_PRIVATE_KEY")
	if !exists {
		return nil, errNoVoterPrivateKey
	}

	choice, exists := os.LookupEnv("PROPOSAL_CHOICE")
	if !exists {
		return nil, errNoChoice
	}

	space, exists := os.LookupEnv("PROPOSAL_SPACE")
	if !exists {
		return nil, errNoSpace
	}

	title, exists := os.LookupEnv("PROPOSAL_TITLE")
	if !exists {
		return nil, errNoProposalTitle
	}

	return &Config{
		Voter: Voter{
			Address:    address,
			PrivateKey: privateKey,
		},
		Proposal: Proposal{
			Choice: choice,
			Space:  space,
			Title:  title,
		},
	}, nil
}
