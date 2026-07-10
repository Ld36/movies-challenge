package usecases

import (
	"context"
	"errors"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/ports"
)

type CreateMovieUseCase struct {
	repository ports.MovieRepository
	publisher  ports.MovieEventPublisher
}

func NewCreateMovieUseCase(repository ports.MovieRepository, publisher ports.MovieEventPublisher) *CreateMovieUseCase {
	return &CreateMovieUseCase{
		repository: repository,
		publisher:  publisher,
	}
}

func (uc *CreateMovieUseCase) Execute(ctx context.Context, id int64, title string, year string) (*domain.Movie, error) {
	movie, err := domain.NewMovie(id, title, year)
	if err != nil {
		return nil, err
	}

	existing, err := uc.repository.FindByID(ctx, movie.ID)
	if err != nil && !errors.Is(err, apperrors.ErrMovieNotFound) {
		return nil, err
	}

	if existing != nil {
		return nil, apperrors.ErrMovieAlreadyExists
	}

	if err := uc.publisher.PublishMovieCreated(ctx, *movie); err != nil {
		return nil, err
	}

	return movie, nil
}
