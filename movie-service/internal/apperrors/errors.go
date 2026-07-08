package apperrors

import "errors"

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrMovieAlreadyExists = errors.New("movie already exists")
	ErrInvalidMovieData   = errors.New("invalid movie data")
)
