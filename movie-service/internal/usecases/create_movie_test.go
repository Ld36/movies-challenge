package usecases

import (
	"context"
	"testing"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/luizdavid/movies-challenge/movie-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMovieUseCase_Execute_Success(t *testing.T) {
	repository := new(mocks.MovieRepositoryMock)
	useCase := NewCreateMovieUseCase(repository)

	ctx := context.Background()

	expectedMovie := &domain.Movie{
		ID:    999,
		Title: "Test Movie",
		Year:  "2026",
	}

	repository.
		On("Create", mock.Anything, mock.MatchedBy(func(movie domain.Movie) bool {
			return movie.ID == expectedMovie.ID &&
				movie.Title == expectedMovie.Title &&
				movie.Year == expectedMovie.Year
		})).
		Return(expectedMovie, nil)

	result, err := useCase.Execute(ctx, 999, "Test Movie", "2026")

	assert.NoError(t, err)
	assert.Equal(t, expectedMovie, result)

	repository.AssertExpectations(t)
}

func TestCreateMovieUseCase_Execute_InvalidMovieData(t *testing.T) {
	repository := new(mocks.MovieRepositoryMock)
	useCase := NewCreateMovieUseCase(repository)

	ctx := context.Background()

	result, err := useCase.Execute(ctx, 0, "", "")

	assert.ErrorIs(t, err, apperrors.ErrInvalidMovieData)
	assert.Nil(t, result)

	repository.AssertNotCalled(t, "Create")
}
