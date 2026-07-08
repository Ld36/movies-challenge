package application

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type MovieService struct {
	repository ports.MovieRepository
}

func NewMovieService(repository ports.MovieRepository) *MovieService {
	return &MovieService{
		repository: repository,
	}
}

func (s *MovieService) GetMovies(ctx context.Context) ([]domain.Movie, error) {
	return s.repository.FindAll(ctx)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int64) (*domain.Movie, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *MovieService) CreateMovie(ctx context.Context, id int64, title string, year string) (*domain.Movie, error) {
	movie, err := domain.NewMovie(id, title, year)
	if err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, *movie)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
