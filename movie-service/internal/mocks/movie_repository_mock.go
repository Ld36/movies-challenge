package mocks

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MovieRepositoryMock struct {
	mock.Mock
}

func (m *MovieRepositoryMock) FindAll(ctx context.Context, page int64, limit int64) ([]domain.Movie, error) {
	args := m.Called(ctx, page, limit)

	movies, _ := args.Get(0).([]domain.Movie)

	return movies, args.Error(1)
}

func (m *MovieRepositoryMock) FindByID(ctx context.Context, id int64) (*domain.Movie, error) {
	args := m.Called(ctx, id)

	movie, _ := args.Get(0).(*domain.Movie)

	return movie, args.Error(1)
}

func (m *MovieRepositoryMock) Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error) {
	args := m.Called(ctx, movie)

	createdMovie, _ := args.Get(0).(*domain.Movie)

	return createdMovie, args.Error(1)
}

func (m *MovieRepositoryMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
