package usecases

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type DeleteMovieUseCase struct {
	repository ports.MovieRepository
}

func NewDeleteMovieUseCase(repository ports.MovieRepository) *DeleteMovieUseCase {
	return &DeleteMovieUseCase{repository: repository}
}

func (uc *DeleteMovieUseCase) Execute(ctx context.Context, id int64) error {
	return uc.repository.Delete(ctx, id)
}
