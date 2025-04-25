package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type tmdbMovie struct {
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	VoteAverage      float64 `json:"vote_average"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	OriginalLanguage string  `json:"original_language"`
}

type tmdbResponse struct {
	Results []tmdbMovie `json:"results"`
}

func fetchMovieFromTMDB(title, langCode, year string) (map[string]interface{}, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	baseURL := "https://api.themoviedb.org/3/search/movie"

	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("query", title)
	if langCode != "" {
		params.Add("with_original_language", langCode)
	}
	if year != "" {
		params.Add("primary_release_year", year)
	}

	finalURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmdbResp tmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, err
	}

	if len(tmdbResp.Results) == 0 {
		return nil, nil
	}

	m := tmdbResp.Results[0]
	if m.Title == "" || m.Overview == "" {
		return nil, nil
	}

	return map[string]interface{}{
		"movie_name":   m.Title,
		"image_url":    "https://image.tmdb.org/t/p/w500" + m.PosterPath,
		"rating":       fmt.Sprintf("%.1f", m.VoteAverage),
		"plot":         m.Overview,
		"release_date": m.ReleaseDate,
		"language":     m.OriginalLanguage,
	}, nil
}
