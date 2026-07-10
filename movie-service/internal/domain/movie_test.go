package domain

import (
	"testing"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestNewMovie_Success(t *testing.T) {
	movie, err := NewMovie(1, " Test Movie ", " 2026 ")

	assert.NoError(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, int64(1), movie.ID)
	assert.Equal(t, "Test Movie", movie.Title)
	assert.Equal(t, "2026", movie.Year)
}

func TestNewMovie_InvalidID(t *testing.T) {
	movie, err := NewMovie(0, "Test Movie", "2026")

	assert.ErrorIs(t, err, apperrors.ErrInvalidMovieData)
	assert.Nil(t, movie)
}

func TestNewMovie_EmptyTitle(t *testing.T) {
	movie, err := NewMovie(1, "   ", "2026")

	assert.ErrorIs(t, err, apperrors.ErrInvalidMovieData)
	assert.Nil(t, movie)
}

func TestNewMovie_EmptyYear(t *testing.T) {
	movie, err := NewMovie(1, "Test Movie", "   ")

	assert.ErrorIs(t, err, apperrors.ErrInvalidMovieData)
	assert.Nil(t, movie)
}
