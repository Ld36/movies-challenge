package usecases

import "github.com/luizdavid/movies-challenge/movie-service/internal/ports"

type MovieUseCases struct {
	ListMovies  *ListMoviesUseCase
	GetMovie    *GetMovieUseCase
	CreateMovie *CreateMovieUseCase
	DeleteMovie *DeleteMovieUseCase
}

func NewMovieUseCases(repository ports.MovieRepository, publisher ports.MovieEventPublisher) *MovieUseCases {
	return &MovieUseCases{
		ListMovies:  NewListMoviesUseCase(repository),
		GetMovie:    NewGetMovieUseCase(repository),
		CreateMovie: NewCreateMovieUseCase(repository, publisher),
		DeleteMovie: NewDeleteMovieUseCase(publisher),
	}
}
