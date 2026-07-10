package ports

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
)

type MovieEventPublisher interface {
	PublishMovieCreated(ctx context.Context, movie domain.Movie) error
	PublishMovieDeleted(ctx context.Context, id int64) error
}
