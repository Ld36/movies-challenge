package grpc

import (
	"errors"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatusError(err error) error {
	switch {
	case errors.Is(err, apperrors.ErrMovieNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, apperrors.ErrMovieAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, apperrors.ErrInvalidMovieData):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
