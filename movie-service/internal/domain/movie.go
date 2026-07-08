package domain

import (
	"strings"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
)

type Movie struct {
	ID    int64  `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Year  string `json:"year" bson:"year"`
}

func NewMovie(id int64, title string, year string) (*Movie, error) {
	movie := &Movie{
		ID:    id,
		Title: strings.TrimSpace(title),
		Year:  strings.TrimSpace(year),
	}

	if err := movie.Validate(); err != nil {
		return nil, err
	}

	return movie, nil
}

func (m *Movie) Validate() error {
	if m.ID <= 0 {
		return apperrors.ErrInvalidMovieData
	}

	if strings.TrimSpace(m.Title) == "" {
		return apperrors.ErrInvalidMovieData
	}

	if strings.TrimSpace(m.Year) == "" {
		return apperrors.ErrInvalidMovieData
	}

	return nil
}
