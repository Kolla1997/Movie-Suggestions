package handler

import (
	"fmt"
	"net/http"
	"time"

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
	if req.TimeOfDay == "" {
		req.TimeOfDay = "any"
	} else if req.TimeOfDay == "Yes" || req.TimeOfDay == "yes" {
		currentHour := time.Now().Hour()
		switch {
		case currentHour >= 5 && currentHour < 12:
			req.TimeOfDay = "Morning"
		case currentHour >= 12 && currentHour < 17:
			req.TimeOfDay = "Afternoon"
		case currentHour >= 17 && currentHour < 21:
			req.TimeOfDay = "Evening"
		default:
			req.TimeOfDay = "Night"
		}
	}
	page := 1
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	movies, err := service.GetMovieSuggestionsList(req, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}
