package auth

import (
	"fmt"
	"time"
)

const (
	HEADER_NAME      = "Authorization"
	TOKEN_TYPE       = "JWT"
	ALGORITHM        = "HS256"
	minSecertKeySize = 32
)

type jwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecertKeySize {
		return nil, fmt.Errorf("invalid key size")
	}
	return &jwtMaker{secretKey: secretKey}, nil
}

func (maker *jwtMaker) CreateToken(username, role string, duration time.Duration) (string, error) {

	return "", nil
}

func (maker *jwtMaker) VerifyToken(token string) (*Payload, error) {
	return nil, nil
}
