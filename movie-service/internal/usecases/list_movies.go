package usecases

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type ListMoviesUseCase struct {
	repository ports.MovieRepository
}

func NewListMoviesUseCase(repository ports.MovieRepository) *ListMoviesUseCase {
	return &ListMoviesUseCase{repository: repository}
}

func (uc *ListMoviesUseCase) Execute(ctx context.Context) ([]domain.Movie, error) {
	return uc.repository.FindAll(ctx)
}
