package app

import (
	"errors"
	"os"
)

// Voter represents the delegate who will vote
// is fetched from the environment variables
type voter struct {
	address    string
	privateKey string
}

type proposal struct {
	choice string
	space  string
	title  string
}

type config struct {
	voter    voter
	proposal proposal
}

var (
	errNoVoterAddress    = errors.New("VOTER_ADDRESS is not set")
	errNoVoterPrivateKey = errors.New("VOTER_PRIVATE_KEY is not set")
	errNoChoice          = errors.New("PROPOSAL_CHOICE is not set")
	errNoSpace           = errors.New("PROPOSAL_SPACE is not set")
	errNoProposalTitle   = errors.New("PROPOSAL_TITLE is not set")
)

func loadConfig() (*config, error) {
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

	return &config{
		voter: voter{
			address:    address,
			privateKey: privateKey,
		},
		proposal: proposal{
			choice: choice,
			space:  space,
			title:  title,
		},
	}, nil
}
