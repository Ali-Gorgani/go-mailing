package models

import (
	"database/sql"
	"fmt"
	"go-mailing/internal/app/utils"
	"strings"
	"time"
)

type CreateUserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreateUserResponse struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashed_password"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserService struct {
	DB *sql.DB
}

func (service *UserService) CreateUser(arg CreateUserParam) (CreateUserResponse, error) {
	arg.Email = strings.ToLower(arg.Email)

	hashedPassword, err := utils.HashPassword(arg.Password)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("create user: %w", err)
	}
	var id int
	row := service.DB.QueryRow(`
		INSERT INTO users (username, hashed_password, email) 
		VALUES ($1, $2, $3) RETURNING id;`, arg.Username, hashedPassword, arg.Email)
	err = row.Scan(&id)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("create user: %w", err)
	}

	// Create the response
	response := CreateUserResponse{
		ID:             id,
		Username:       arg.Username,
		HashedPassword: hashedPassword,
		Email:          arg.Email,
		CreatedAt:      time.Now(),
	}

	return response, nil
}

type SignInParam struct {
	Username string `query:"username"`
	Password string `query:"password"`
}

func (service *UserService) SignIn(args SignInParam) (CreateUserResponse, error) {
	var user CreateUserResponse
	row := service.DB.QueryRow(`
		SELECT id, username, hashed_password, email, created_at 
		FROM users 
		WHERE username = $1;`, args.Username)
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email, &user.CreatedAt)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("sign in: %w", err)
	}

	err = utils.ComparePassword(user.HashedPassword, args.Password)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("sign in: %w", err)
	}

	return user, nil
}

func (service *UserService) Update() {
	// ...
}

func (service *UserService) Delete() {
	// ...
}
