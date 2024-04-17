package core

import (
	"github.com/gofiber/fiber/v3"
)

func InitCoreRoutes(fiber *fiber.App) {
	routes := fiber.Group("/")
	routes.Get("/health", HealthCheck)
}
