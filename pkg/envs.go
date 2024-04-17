package pkg

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/dozto/dozto-aster-auth/pkg/helper"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type MyConfig struct {
	Port      string `env:"PORT" envDefault:"3000"`
	MongoUri  string `env:"MONGODB_URI" validate:"required,startswith=mongodb"`
	MongoName string `env:"MONGODB_NAME" validate:"required"`
}

func LoadEnvs() (config MyConfig) {

	STAGE := os.Getenv("GO_ENV")
	if STAGE == "" {
		STAGE = "dev"
	}

	// load .env file to env for non prod
	if STAGE != "prod" {
		log.Info().Msg("loading .env file config...")
		root := helper.RootPath()
		envPath := root + "/.env." + STAGE

		if err := godotenv.Load(envPath); err != nil {
			log.Fatal().Err(err).Msgf("fail load env file, required one for stage: %s", STAGE)
		}
	}

	// parse env to config

	if err := env.Parse(&config); err != nil {
		log.Fatal().Err(err).Msgf("fail parse envs to config, check system env or .env file")
	}

	// Validate config
	if err := Valtor.Validate(config); err != nil {
		log.Fatal().Err(err).Msg("fail validate env varialbes")
	}

	log.Info().Msgf("success load config for stage: %s", STAGE)
	return config
}
