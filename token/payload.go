package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	ID uuid.UUID `json:""`
	Username string `json:"username"`
	IssueAt time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId,err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: tokenId,
		Username: username,
		IssueAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload,nil
}

// Valid check if token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpired
	}
	return nil
}