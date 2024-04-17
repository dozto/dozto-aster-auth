package app

import (
	"github.com/dozto/dozto-aster-auth/app/auth"
	"github.com/dozto/dozto-aster-auth/app/core"
	"github.com/dozto/dozto-aster-auth/app/user"
	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitFiberServer(config *pkg.MyConfig, db *mongo.Database) *fiber.App {
	userModel := user.NewUserModel(db, "users")

	fiber := fiber.New()
	core.InitCoreRoutes(fiber)
	auth.InitAuthRoutes(fiber, userModel)
	user.InitUserRoutes(fiber, userModel)

	return fiber
}
