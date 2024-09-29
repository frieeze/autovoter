package signer

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type Vote struct {
	From      string `json:"from"`
	Space     string `json:"space"`
	Timestamp int64  `json:"timestamp"`
	Choice    int    `json:"choice"`
	Proposal  string `json:"proposal"`
}

type Signer struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

func New(address, privateKey string) (*Signer, error) {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return &Signer{
		Address:    address,
		PrivateKey: pk,
	}, nil
}

const message = `{
    "from": "0xc1c39b466a3660e64bfc5c256e6b8e7083957a4a",
    "space": "cvx.eth",
    "timestamp": "1727621658",
    "proposal": "0x9b1faf762db03057ec16ad3b347548a4d19bbe35d01c8b20cd725239e8c89028",
    "choice": "{\"477\":3}",
    "reason": "",
    "app": "snapshot",
    "metadata": "{}"
}`

func (s *Signer) Vote(choice int, proposal, space string) (*Vote, string, error) {
	vote := &Vote{
		From:      s.Address,
		Space:     space,
		Timestamp: time.Now().Unix(),
		Choice:    choice,
		Proposal:  proposal,
	}
	sig, err := s.signVote(vote)
	return vote, sig, err
}

func (s *Signer) signVote(vote *Vote) (string, error) {
	data := buildMessage(vote)
	hash := data.TypeHash("Vote")

	js, err := json.Marshal(data.Map())
	if err != nil {
		return "", fmt.Errorf("failed to marshal vote: %w", err)
	}
	fmt.Printf("signed data:\n %s \n", string(js))

	signature, err := crypto.Sign(hash, s.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign vote: %w", err)
	}
	fmt.Printf("signature:\n %s \n", hexutil.Encode(signature))
	return hexutil.Encode(signature), nil
}

func buildMessage(vote *Vote) apitypes.TypedData {
	return apitypes.TypedData{
		PrimaryType: "Vote",
		Domain: apitypes.TypedDataDomain{
			Name:    "snapshot",
			Version: "0.1.4",
		},
		Types: apitypes.Types{
			"Vote": []apitypes.Type{
				{Name: "from", Type: "address"},
				{Name: "space", Type: "string"},
				{Name: "timestamp", Type: "uint64"},
				{Name: "proposal", Type: "bytes32"},
				{Name: "choice", Type: "string"},
				{Name: "reason", Type: "string"},
				{Name: "app", Type: "string"},
				{Name: "metadata", Type: "string"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
			},
		},
		Message: apitypes.TypedDataMessage{
			"from":      vote.From,
			"space":     vote.Space,
			"timestamp": vote.Timestamp,
			"proposal":  vote.Proposal,
			"choice":    fmt.Sprintf("{\"%d\":1}", vote.Choice),
			"reason":    "",
			"app":       "flying-penguin",
			"metadata":  "{}",
		},
	}
}
