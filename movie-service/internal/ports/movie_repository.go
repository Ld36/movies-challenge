package ports

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
)

type MovieRepository interface {
	FindAll(ctx context.Context) ([]domain.Movie, error)
	FindByID(ctx context.Context, id int64) (*domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id int64) error
}
