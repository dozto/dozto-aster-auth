package mongodb

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string, dbName string, timeout time.Duration) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Info().Msg("connecting to mongodb...")

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	log.Info().Msgf("success connect to mongodb: %s", dbName)
	return client.Database(dbName)
}

func CloseClient(db *mongo.Database) error {
	return db.Client().Disconnect(context.Background())
}

func CreateIndex(coll *mongo.Collection, indexModel []mongo.IndexModel) ([]string, error) {
	names, err := coll.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create index")
	}

	log.Info().Msgf("Collection %s Index Created: %s ", coll.Name(), strings.Join(names, ";"))

	return names, nil
}
