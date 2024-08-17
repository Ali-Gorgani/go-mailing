package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID
	Username  string
	Role      string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(username, role string, Duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token id: %w", err)
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(Duration),
	}
	return payload, nil
}
