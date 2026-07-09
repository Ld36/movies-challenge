package routes

import (
	"github.com/gin-gonic/gin"
	httpadapter "github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/http"
)

func RegisterRoutes(router *gin.Engine, movieHandler *httpadapter.MovieHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	movies := router.Group("/movies")
	{
		movies.GET("", movieHandler.GetMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.POST("", movieHandler.CreateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)
	}
}
