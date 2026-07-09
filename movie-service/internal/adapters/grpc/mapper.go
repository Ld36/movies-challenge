package grpc

import (
	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	moviepb "github.com/luizdavid/movies-challenge/movie-service/proto"
)

func ToProtoMovie(movie domain.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:    movie.ID,
		Title: movie.Title,
		Year:  movie.Year,
	}
}

func ToProtoMovies(movies []domain.Movie) []*moviepb.Movie {
	result := make([]*moviepb.Movie, 0, len(movies))

	for _, movie := range movies {
		result = append(result, ToProtoMovie(movie))
	}

	return result
}
