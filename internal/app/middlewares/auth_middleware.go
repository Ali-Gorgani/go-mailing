package middlewares

import (
	"go-mailing/internal/app/models"
	"go-mailing/pkg/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type AuthMiddleware struct {
	SessionService *models.SessionService
	TokenMaker     auth.Maker
}

func (amw *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorizationHeader := ctx.Request().Header.Get(authorizationHeaderKey)
		if authorizationHeader == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unsupported authorization type"})
		}

		accessToken := fields[1]
		payload, err := amw.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired access token"})
		}

		ctx.Set(authorizationPayloadKey, payload)
		return next(ctx)
	}
}
