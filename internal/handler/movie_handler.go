package handler

import (
	"net/http"

	"github.com/dev1force/Movie-Suggestions.git/internal/model"
	"github.com/dev1force/Movie-Suggestions.git/internal/service"
	"github.com/gin-gonic/gin"
)

func GetMovieSuggestionList(c *gin.Context) {
	var req model.MovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	movies, err := service.GetMovieSuggestionsList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}
