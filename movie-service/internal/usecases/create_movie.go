package usecases

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type CreateMovieUseCase struct {
	repository ports.MovieRepository
}

func NewCreateMovieUseCase(repository ports.MovieRepository) *CreateMovieUseCase {
	return &CreateMovieUseCase{repository: repository}
}

func (uc *CreateMovieUseCase) Execute(ctx context.Context, id int64, title string, year string) (*domain.Movie, error) {
	movie, err := domain.NewMovie(id, title, year)
	if err != nil {
		return nil, err
	}

	return uc.repository.Create(ctx, *movie)
}
