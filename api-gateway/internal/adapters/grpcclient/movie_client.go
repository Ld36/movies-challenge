package grpcclient

import (
	"context"

	moviepb "github.com/luizdavid/movies-challenge/api-gateway/proto"
	"google.golang.org/grpc"
)

type MovieClient struct {
	client moviepb.MovieServiceClient
}

func NewMovieClient(conn *grpc.ClientConn) *MovieClient {
	return &MovieClient{
		client: moviepb.NewMovieServiceClient(conn),
	}
}

func (c *MovieClient) GetMovies(ctx context.Context) (*moviepb.GetMoviesResponse, error) {
	return c.client.GetMovies(ctx, &moviepb.GetMoviesRequest{})
}

func (c *MovieClient) GetMovieByID(ctx context.Context, id int64) (*moviepb.GetMovieByIdResponse, error) {
	return c.client.GetMovieById(ctx, &moviepb.GetMovieByIdRequest{
		Id: id,
	})
}

func (c *MovieClient) CreateMovie(ctx context.Context, id int64, title string, year string) (*moviepb.CreateMovieResponse, error) {
	return c.client.CreateMovie(ctx, &moviepb.CreateMovieRequest{
		Id:    id,
		Title: title,
		Year:  year,
	})
}

func (c *MovieClient) DeleteMovie(ctx context.Context, id int64) (*moviepb.DeleteMovieResponse, error) {
	return c.client.DeleteMovie(ctx, &moviepb.DeleteMovieRequest{
		Id: id,
	})
}
