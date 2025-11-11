package routes

import (
	"backend/web/controller"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	LogMiddleware     fiber.Handler
	BranchsController *controller.BranchsController
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.LogMiddleware)
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	branch := c.App.Group("branch")
	branch.Post("/management", c.BranchsController.AddNewManagement)
	branch.Post("/", c.BranchsController.AddNewBranch)
}

func (c *RouteConfig) SetupAuthRoute() {}
