package usecases

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type GetMovieUseCase struct {
	repository ports.MovieRepository
}

func NewGetMovieUseCase(repository ports.MovieRepository) *GetMovieUseCase {
	return &GetMovieUseCase{repository: repository}
}

func (uc *GetMovieUseCase) Execute(ctx context.Context, id int64) (*domain.Movie, error) {
	return uc.repository.FindByID(ctx, id)
}
