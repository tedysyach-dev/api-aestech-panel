package middlewares

import (
	"backend/core/utils"
	"backend/web/model"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(tokenUtil *utils.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.ErrUnauthorized
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		request := &model.VerifyUserRequest{Token: token}

		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
