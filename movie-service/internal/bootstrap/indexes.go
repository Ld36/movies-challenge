package bootstrap

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func EnsureMovieIndexes(ctx context.Context, collection *mongo.Collection, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
		},
		Options: options.Index().
			SetName("idx_movies_id_unique").
			SetUnique(true),
	}

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	logger.Info("movie indexes ensured", zap.String("index", indexName))

	return nil
}
