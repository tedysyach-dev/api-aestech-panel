package routes

import "github.com/gofiber/fiber/v2"

type RouteConfig struct {
	App           *fiber.App
	LogMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.LogMiddleware)
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {}

func (c *RouteConfig) SetupAuthRoute() {}
