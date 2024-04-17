package test

import (
	"time"

	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/dozto/dozto-aster-auth/pkg/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init() *mongo.Database {
	// Init Global Logger
	pkg.InitPrettyLogger()

	// LoadEnvs
	config := pkg.LoadEnvs()
	log.Info().Msg("test stage init complete")

	// Connect to MongoDB
	db := mongodb.Connect(config.MongoUri, config.MongoName, 20*time.Second)

	return db
}

func CleanUp() {
	log.Info().Msg("test stage cleanup complete")
}
