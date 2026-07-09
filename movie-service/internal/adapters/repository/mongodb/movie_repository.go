package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(collection *mongo.Collection) *MovieRepository {
	return &MovieRepository{
		collection: collection,
	}
}

func (r *MovieRepository) FindAll(ctx context.Context, page int64, limit int64) ([]domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "id", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []domain.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) FindByID(ctx context.Context, id int64) (*domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var movie domain.Movie

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&movie)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.ErrMovieNotFound
	}

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := r.FindByID(ctx, movie.ID)
	if err != nil && !errors.Is(err, apperrors.ErrMovieNotFound) {
		return nil, err
	}

	if existing != nil {
		return nil, apperrors.ErrMovieAlreadyExists
	}

	if _, err := r.collection.InsertOne(ctx, movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperrors.ErrMovieNotFound
	}

	return nil
}
