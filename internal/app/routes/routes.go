package routes

import (
	"database/sql"
	"go-mailing/internal/app/controllers"
	"go-mailing/internal/app/models"

	"github.com/labstack/echo/v4"
)

func NewRouter(db *sql.DB) *echo.Echo {

	// Set up the services
	userService := &models.UserService{
		DB: db,
	}

	// Set up the controllers
	userC := &controllers.UserController{
		UserService: userService,
	}

	smsC := &controllers.SMSController{}

	// Set up the routes
	e := echo.New()
	e.POST("/signup", userC.SignUp)
	e.POST("/signin", userC.SignIn)
	e.POST("/signout", userC.SignOut)

	userGroup := e.Group("/user")
	userGroup.GET("/give-providers", smsC.GiveProviders)
	userGroup.POST("/post-sms", smsC.PostSMS)

	e.Static("/swagger", "swagger")
	return e
}
