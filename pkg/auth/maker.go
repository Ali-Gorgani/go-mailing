package auth

import "time"

type Maker interface {
	CreateToken(username, role string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
