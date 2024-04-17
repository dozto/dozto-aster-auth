package main

import (
	"time"

	"github.com/dozto/dozto-aster-auth/app"
	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/dozto/dozto-aster-auth/pkg/mongodb"
	"github.com/rs/zerolog/log"
)

func main() {
	// Init Global Logger
	pkg.InitPrettyLogger()

	// LoadEnvs
	myConfig := pkg.LoadEnvs()
	log.Info().Any("envs", myConfig.Port)

	// Connect to MongoDB
	db := mongodb.Connect(myConfig.MongoUri, myConfig.MongoName, 20*time.Second)

	// Start App Server
	server := app.InitFiberServer(&myConfig, db)
	log.Info().Msgf("server is running on port: %s", myConfig.Port)
	server.Listen(":" + myConfig.Port)
}
