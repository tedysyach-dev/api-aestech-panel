package controller

import (
	"backend/core/utils"
	"backend/web/model"
	"backend/web/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BranchsController struct {
	Log     *logrus.Logger
	Service *service.BranchsService
}

func NewBrandsController(service *service.BranchsService, logger *logrus.Logger) *BranchsController {
	return &BranchsController{
		Log:     logger,
		Service: service,
	}
}

func (c *BranchsController) AddNewManagement(ctx *fiber.Ctx) error {
	request := new(model.CreateManagementRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	res, err := c.Service.AddNewManagement(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create management : %+v", err)
		return err
	}

	return ctx.JSON(utils.WebResponse[*model.CreateManagementResponse]{Status: true, Message: "Success", Code: fiber.StatusOK, Data: res})
}

func (c *BranchsController) AddNewBranch(ctx *fiber.Ctx) error {
	request := new(model.CreateBranchRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	res, err := c.Service.AddNewBranch(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create management : %+v", err)
		return err
	}

	return ctx.JSON(utils.WebResponse[*model.CreateBranchResponse]{Status: true, Message: "Success", Code: fiber.StatusOK, Data: res})
}
