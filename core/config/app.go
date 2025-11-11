package config

import (
	"backend/core/middlewares"
	"backend/core/routes"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	appMode := config.Config.GetString("app.development")

	if appMode == "dev" {
		config.App.Use(cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
			ExposeHeaders: "Content-Length",
			// AllowCredentials: true,
		}))
	} else {
		config.App.Use(cors.New(cors.Config{
			AllowOrigins:  config.Config.GetString("allowedWeb"),
			AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
			ExposeHeaders: "Content-Length",
			// AllowCredentials: true,
		}))
	}

	logMiddleware := middlewares.RequestLogger(config.Log)

	routeConfig := routes.RouteConfig{
		App:           config.App,
		LogMiddleware: logMiddleware,
	}

	routeConfig.Setup()
}
