package usecases

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type DeleteMovieUseCase struct {
	publisher ports.MovieEventPublisher
}

func NewDeleteMovieUseCase(publisher ports.MovieEventPublisher) *DeleteMovieUseCase {
	return &DeleteMovieUseCase{
		publisher: publisher,
	}
}

func (uc *DeleteMovieUseCase) Execute(ctx context.Context, id int64) error {
	return uc.publisher.PublishMovieDeleted(ctx, id)
}
