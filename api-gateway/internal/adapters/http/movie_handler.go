package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/grpcclient"
	moviepb "github.com/luizdavid/movies-challenge/api-gateway/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieHandler struct {
	client *grpcclient.MovieClient
}

type createMovieRequest struct {
	ID    int64  `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Year  string `json:"year" binding:"required"`
}

func NewMovieHandler(client *grpcclient.MovieClient) *MovieHandler {
	return &MovieHandler{
		client: client,
	}
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	response, err := h.client.GetMovies(c.Request.Context())
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response.Movies,
	})
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid movie id",
		})
		return
	}

	response, err := h.client.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response.Movie,
	})
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var request createMovieRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	response, err := h.client.CreateMovie(
		c.Request.Context(),
		request.ID,
		request.Title,
		request.Year,
	)

	if err != nil {
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": response.Movie,
	})
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid movie id",
		})
		return
	}

	_, err = h.client.DeleteMovie(c.Request.Context(), id)
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": true,
	})
}

func parseIDParam(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	return strconv.ParseInt(idParam, 10, 64)
}

func handleGRPCError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	st, ok := status.FromError(err)
	if !ok {
		c.JSON(statusCode, gin.H{"error": message})
		return
	}

	message = st.Message()

	switch st.Code() {
	case codes.NotFound:
		statusCode = http.StatusNotFound

	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest

	case codes.AlreadyExists:
		statusCode = http.StatusConflict

	default:
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

var _ = errors.New
var _ = moviepb.Movie{}
