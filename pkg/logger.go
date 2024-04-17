package pkg

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitPrettyLogger() {
	stage := os.Getenv("STAGE")
	if stage != "prod" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
