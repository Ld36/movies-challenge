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

func (uc *ListMoviesUseCase) Execute(ctx context.Context, page int64, limit int64) ([]domain.Movie, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return uc.repository.FindAll(ctx, page, limit)
}
