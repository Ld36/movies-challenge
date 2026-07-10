package mocks

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MovieEventPublisherMock struct {
	mock.Mock
}

func (m *MovieEventPublisherMock) PublishMovieCreated(ctx context.Context, movie domain.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MovieEventPublisherMock) PublishMovieDeleted(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
