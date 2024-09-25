package signer

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

const message = `{"from":"%s","space":"%s","timestamp":%d,"proposal":"%s","choice":"{\"%d\":1}","reason":"","app":"flying-penguin","metadata":"{}"}`

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
	msg := fmt.Sprintf(message, vote.From, vote.Space, vote.Timestamp, vote.Proposal, vote.Choice)
	hash := crypto.Keccak256([]byte(msg))
	signature, err := crypto.Sign(hash, s.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign vote: %w", err)
	}
	return hexutil.Encode(signature), nil
}
