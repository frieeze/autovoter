package snapshot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/frieeze/autovoter/http"
)

// Client is a snapshot client
type Client struct {
	HUB       string
	Sequencer string
}

var (
	// ErrNoActiveProposal is returned when no active proposal is found
	ErrNoActiveProposal = errors.New("no active proposal found")
	// ErrNoMatchingProposal is returned when no proposal is found
	ErrNoMatchingProposal = errors.New("no matching proposal found")
	// ErrNoMatchingChoice is returned when no choice is found
	ErrNoMatchingChoice = errors.New("no matching choice found")
)

type ssProposal struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Choices []string `json:"choices"`
}

func (c *Client) GetProposal(ctx context.Context, space, title, label string) (string, int, error) {
	proposals, err := c.queryProposals(ctx, space, title)
	if err != nil {
		return "", 0, err
	}

	if len(proposals) == 0 {
		return "", 0, ErrNoMatchingProposal
	}

	var (
		pId     string
		choices []string
	)
	for _, p := range proposals {
		if strings.HasPrefix(p.Title, title) {
			pId = p.ID
			choices = p.Choices
			break
		}
	}

	if pId == "" {
		return "", 0, ErrNoMatchingProposal
	}

	for idx, choice := range choices {
		if strings.Contains(choice, label) {
			return pId, idx, nil
		}
	}
	return pId, 0, ErrNoMatchingChoice
}

func (c *Client) queryProposals(ctx context.Context, space, title string) ([]ssProposal, error) {
	type response struct {
		Data struct {
			Proposals []ssProposal `json:"proposals"`
		} `json:"data"`
	}

	var (
		resp  = &response{}
		query = fmt.Sprintf(`
		query Proposals {
			proposals(
				where: {
					space_in: %s,
					title_contains: "%s",
					state: "active"
				},
				orderBy: "created",
				orderDirection: desc
			) {
				id
				title
				choices
			}
		}`, space, title)
	)

	err := http.Post(ctx, c.HUB, queryToBody(query), resp)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch proposal: %w", err)
	}

	return resp.Data.Proposals, nil
}

func queryToBody(query string) map[string]string {
	return map[string]string{"query": query}
}

const voteBody = `{"address":"%s","sig":"%s","data":{"domain":{"name":"snapshot","version":"0.1.4"},"types":{"Vote":[{"name":"from","type":"address"},{"name":"space","type":"string"},{"name":"timestamp","type":"uint64"},{"name":"proposal","type":"bytes32"},{"name":"choice","type":"string"},{"name":"reason","type":"string"},{"name":"app","type":"string"},{"name":"metadata","type":"string"}]},"message":{"from":"%s","space":"%s","timestamp":%d,"proposal":"%s","choice":"{\"%d\":1}","reason":"","app":"flying-penguin","metadata":"{}"}}}`

func (c *Client) SendVote(ctx context.Context, message *autovoter.VoteMessage, sig string) error {
	body := fmt.Sprintf(voteBody,
		message.From,
		sig,
		message.From,
		message.Space,
		message.Timestamp,
		message.Proposal,
		message.Choice,
	)
	return http.Post(ctx, c.HUB, body, nil)
}
