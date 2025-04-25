package main

import (
	"log"
	"os"

	"github.com/dev1force/Movie-Suggestions.git/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.POST("/movie-suggestions", handler.GetMovieSuggestionList)

	port := os.Getenv("PORT")
	log.Printf("Running at http://localhost:%s", port)
	r.Run(":" + port)
}
