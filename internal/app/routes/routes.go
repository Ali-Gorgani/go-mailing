package routes

import (
	"database/sql"
	"go-mailing/configs"
	"go-mailing/internal/app/controllers"
	"go-mailing/internal/app/middlewares"
	"go-mailing/internal/app/models"
	"go-mailing/pkg/auth"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func NewRouter(db *sql.DB, log *logrus.Logger, cfg configs.Config) *echo.Echo {

	jwtMaker, err := auth.NewJwtMaker(cfg.AccessToken.SecretKey)
	if err != nil {
		log.Fatalf("new router: failed to create jwt maker: %v", err)
	}

	// Set up the services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	// Set up the controllers
	userC := &controllers.UserController{
		UserService:    userService,
		SessionService: sessionService,
		TokenMaker:     jwtMaker,
		Log:            log,
		Config:         cfg,
	}

	smsC := &controllers.SMSController{}

	// Set up the middleware
	authMiddleware := &middlewares.AuthMiddleware{
		SessionService: sessionService,
		TokenMaker:     jwtMaker,
	}

	// Set up the routes
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.POST("/signup", userC.SignUp)
	e.POST("/signin", userC.SignIn)

	userGroup := e.Group("/user")
	userGroup.Use(authMiddleware.Handle)
	userGroup.POST("/signout", userC.SignOut)
	userGroup.GET("/:username", userC.UserInfo)
	userGroup.GET("/give-providers", smsC.GiveProviders)
	userGroup.POST("/post-sms", smsC.PostSMS)

	e.Static("/swagger", "swagger")
	return e
}
