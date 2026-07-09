package grpc

import (
	"context"

	"github.com/luizdavid/movies-challenge/movie-service/internal/usecases"
	moviepb "github.com/luizdavid/movies-challenge/movie-service/proto"
)

type MovieHandler struct {
	moviepb.UnimplementedMovieServiceServer
	useCases *usecases.MovieUseCases
}

func NewMovieHandler(useCases *usecases.MovieUseCases) *MovieHandler {
	return &MovieHandler{
		useCases: useCases,
	}
}

func (h *MovieHandler) GetMovies(ctx context.Context, req *moviepb.GetMoviesRequest) (*moviepb.GetMoviesResponse, error) {
	movies, err := h.useCases.ListMovies.Execute(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, ToStatusError(err)
	}

	return &moviepb.GetMoviesResponse{
		Movies: ToProtoMovies(movies),
	}, nil
}

func (h *MovieHandler) GetMovieById(ctx context.Context, req *moviepb.GetMovieByIdRequest) (*moviepb.GetMovieByIdResponse, error) {
	movie, err := h.useCases.GetMovie.Execute(ctx, req.Id)
	if err != nil {
		return nil, ToStatusError(err)
	}

	return &moviepb.GetMovieByIdResponse{
		Movie: ToProtoMovie(*movie),
	}, nil
}

func (h *MovieHandler) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {
	movie, err := h.useCases.CreateMovie.Execute(ctx, req.Id, req.Title, req.Year)
	if err != nil {
		return nil, ToStatusError(err)
	}

	return &moviepb.CreateMovieResponse{
		Movie: ToProtoMovie(*movie),
	}, nil
}

func (h *MovieHandler) DeleteMovie(ctx context.Context, req *moviepb.DeleteMovieRequest) (*moviepb.DeleteMovieResponse, error) {
	if err := h.useCases.DeleteMovie.Execute(ctx, req.Id); err != nil {
		return nil, ToStatusError(err)
	}

	return &moviepb.DeleteMovieResponse{
		Deleted: true,
	}, nil
}
