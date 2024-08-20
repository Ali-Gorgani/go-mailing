package controllers

import (
	"go-mailing/configs"
	"go-mailing/internal/app/models"
	"go-mailing/pkg/auth"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserService    *models.UserService
	SessionService *models.SessionService
	TokenMaker     auth.Maker
	Log            *logrus.Logger
	Config         configs.Config
}

type SignUpRequest struct {
	Username string `query:"username" validate:"required"`
	Password string `query:"password" validate:"required"`
	Email    string `query:"email" validate:"required,email"`
}

type SignUpResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (controller *UserController) SignUp(ctx echo.Context) error {
	request := new(SignUpRequest)
	if err := ctx.Bind(request); err != nil {
		controller.Log.Errorf("SignUp: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid query parameters")
	}

	// TODO: Validate the request

	userParam := models.CreateUserParam{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}
	user, err := controller.UserService.CreateUser(userParam)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			controller.Log.Errorf("SignUp: %v", err)
			return ctx.String(http.StatusConflict, "Username or email already exists")
		}
		controller.Log.Errorf("SignUp: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
	rsp := SignUpResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return ctx.JSON(http.StatusOK, rsp)
}

type SignInRequest struct {
	Username string `query:"username" validate:"required"`
	Password string `query:"password" validate:"required"`
}

type SignInResponse struct {
	SessionID            uuid.UUID      `json:"session_id"`
	AccessToken          string         `json:"access_token"`
	AccessTokenExpiresAt time.Time      `json:"access_token_expires_at"`
	User                 SignUpResponse `json:"user"`
}

func (controller *UserController) SignIn(ctx echo.Context) error {
	request := new(SignInRequest)
	err := ctx.Bind(&request)
	if err != nil {
		controller.Log.Errorf("SignIn: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid query parameters")
	}

	// TODO: Validate the request

	signinParam := models.SignInParam{
		Username: request.Username,
		Password: request.Password,
	}
	user, err := controller.UserService.SignIn(ctx.Request().Context(), signinParam)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			controller.Log.Errorf("SignIn: %v", err)
			return ctx.String(http.StatusUnauthorized, "Invalid username or password.")
		}
		controller.Log.Errorf("SignIn: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// TODO: Implement the session with JWT
	token, payload, err := controller.TokenMaker.CreateToken(user.Username, "user", controller.Config.AccessToken.Duration)
	if err != nil {
		controller.Log.Errorf("SignIn: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	session, err := controller.SessionService.CreateSession(ctx.Request().Context(), models.CreateSessionParams{
		ID:        payload.ID,
		Username:  payload.Username,
		Token:     token,
		IsBlocked: false,
		ExpiresAt: payload.ExpiredAt,
	})
	if err != nil {
		controller.Log.Errorf("SignIn: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	rsp := SignInResponse{
		SessionID:            session.ID,
		AccessToken:          token,
		AccessTokenExpiresAt: payload.ExpiredAt,
		User: SignUpResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}
	return ctx.JSON(http.StatusOK, rsp)
}

func (controller *UserController) SignOut(ctx echo.Context) error {
	// TODO: Destroy the JWT session
	return nil
}

func (controller *UserController) UserInfo(ctx echo.Context) error {
	username := ctx.Param("username")
	user, err := controller.UserService.GetUser(username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			controller.Log.Errorf("UserInfo: %v", err)
			return ctx.String(http.StatusNotFound, "User not found")
		}
		controller.Log.Errorf("UserInfo: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	rsp := SignUpResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return ctx.JSON(http.StatusOK, rsp)
}
