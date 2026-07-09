package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func SeedMovies(ctx context.Context, collection *mongo.Collection, filePath string, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	if count > 0 {
		logger.Info("movies seed skipped, collection already has data", zap.Int64("count", count))
		return nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var movies []domain.Movie

	if err := json.Unmarshal(file, &movies); err != nil {
		return err
	}

	if len(movies) == 0 {
		return errors.New("movies seed file is empty")
	}

	documents := make([]interface{}, 0, len(movies))

	for _, movie := range movies {
		documents = append(documents, movie)
	}

	if _, err := collection.InsertMany(ctx, documents); err != nil {
		return err
	}

	logger.Info("movies seed completed", zap.Int("total", len(movies)))

	return nil
}
