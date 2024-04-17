package user

import (
	"github.com/gofiber/fiber/v3"
)

func InitUserRoutes(fiber *fiber.App, userModel *UserModel) {
	controller := NewUserController(userModel)

	routes := fiber.Group("/users")
	routes.Post("/", controller.Create)
	routes.Get("/:id", controller.GetById)
	// routes.Get("/me", controller.GetCurrentUser)
	// routes.Put("/:id", controller.UpdateById)
	// routes.Delete("/:id", controller.DeleteById)
}
