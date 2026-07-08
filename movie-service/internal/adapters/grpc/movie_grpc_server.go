package grpc

import (
	"context"
	"errors"

	"github.com/luizdavid/movies-challenge/movie-service/internal/apperrors"
	"github.com/luizdavid/movies-challenge/movie-service/internal/application"
	"github.com/luizdavid/movies-challenge/movie-service/internal/domain"
	moviepb "github.com/luizdavid/movies-challenge/movie-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieGRPCServer struct {
	moviepb.UnimplementedMovieServiceServer
	service *application.MovieService
}

func NewMovieGRPCServer(service *application.MovieService) *MovieGRPCServer {
	return &MovieGRPCServer{
		service: service,
	}
}

func (s *MovieGRPCServer) GetMovies(ctx context.Context, req *moviepb.GetMoviesRequest) (*moviepb.GetMoviesResponse, error) {
	movies, err := s.service.GetMovies(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &moviepb.GetMoviesResponse{
		Movies: make([]*moviepb.Movie, 0, len(movies)),
	}

	for _, movie := range movies {
		response.Movies = append(response.Movies, toProtoMovie(movie))
	}

	return response, nil
}

func (s *MovieGRPCServer) GetMovieById(ctx context.Context, req *moviepb.GetMovieByIdRequest) (*moviepb.GetMovieByIdResponse, error) {
	movie, err := s.service.GetMovieByID(ctx, req.Id)
	if errors.Is(err, apperrors.ErrMovieNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moviepb.GetMovieByIdResponse{
		Movie: toProtoMovie(*movie),
	}, nil
}

func (s *MovieGRPCServer) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {
	movie, err := s.service.CreateMovie(ctx, req.Id, req.Title, req.Year)
	if errors.Is(err, apperrors.ErrInvalidMovieData) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, apperrors.ErrMovieAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moviepb.CreateMovieResponse{
		Movie: toProtoMovie(*movie),
	}, nil
}

func (s *MovieGRPCServer) DeleteMovie(ctx context.Context, req *moviepb.DeleteMovieRequest) (*moviepb.DeleteMovieResponse, error) {
	err := s.service.DeleteMovie(ctx, req.Id)
	if errors.Is(err, apperrors.ErrMovieNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moviepb.DeleteMovieResponse{
		Deleted: true,
	}, nil
}

func toProtoMovie(movie domain.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:    movie.ID,
		Title: movie.Title,
		Year:  movie.Year,
	}
}
