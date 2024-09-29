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
	hash, err := data.HashStruct("Vote", data.Message)
	if err != nil {
		return "", fmt.Errorf("failed to hash vote: %w", err)
	}
	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("failed to hash domain: %w", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(hash)))

	signature, err := crypto.Sign(crypto.Keccak256Hash(rawData).Bytes(), s.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign vote: %w", err)
	}
	if signature[64] != 27 && signature[64] != 28 {
		signature[64] += 27
	}

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
			"timestamp": fmt.Sprintf("%d", vote.Timestamp),
			"proposal":  vote.Proposal,
			"choice":    fmt.Sprintf("{\"%d\":1}", vote.Choice),
			"reason":    "",
			"app":       "flying-penguin",
			"metadata":  "{}",
		},
	}
}
