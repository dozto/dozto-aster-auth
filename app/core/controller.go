package core

import "github.com/gofiber/fiber/v3"

func HealthCheck(c fiber.Ctx) error {
	return c.SendString("OK!")
}
