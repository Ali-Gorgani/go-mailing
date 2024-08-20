package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	DB *sql.DB
}

type Session struct {
	ID        uuid.UUID
	Username  string
	Token     string
	IsBlocked bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreateSessionParams struct {
	ID        uuid.UUID
	Username  string
	Token     string
	IsBlocked bool
	ExpiresAt time.Time
}

func (service *SessionService) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := service.DB.QueryRow(`
		INSERT INTO sessions (id, username, token, is_blocked, expires_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING *;`, arg.ID, arg.Username, arg.Token, arg.IsBlocked, arg.ExpiresAt)
	var s Session
	err := row.Scan(&s.ID, &s.Username, &s.Token, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}
	return s, nil
}

func (service *SessionService) GetSessionByToken(ctx context.Context, token string) (Session, error) {
	row := service.DB.QueryRow(`
		SELECT * FROM sessions WHERE token=$1;`, token)
	var s Session
	err := row.Scan(&s.ID, &s.Username, &s.Token, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return Session{}, fmt.Errorf("get session by token: %w", err)
	}
	return s, nil
}
