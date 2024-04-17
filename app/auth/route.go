package auth

import (
	"github.com/dozto/dozto-aster-auth/app/user"
	"github.com/gofiber/fiber/v3"
)

func InitAuthRoutes(fiber *fiber.App, userModel *user.UserModel) {
	controller := NewAuthController(userModel)

	routes := fiber.Group("/auth")
	routes.Post("/login", controller.Login)
}
