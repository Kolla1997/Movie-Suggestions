package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dev1force/Movie-Suggestions.git/internal/model"
	"github.com/sashabaranov/go-openai"
)

func parseFormattedMovies(raw string) []model.ParsedMovie {
	var movies []model.ParsedMovie

	entries := strings.Split(raw, ",")
	re := regexp.MustCompile(`^(.*?) \((\d{4})\) \[(\w{2})\]$`)

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		match := re.FindStringSubmatch(entry)
		if len(match) == 4 {
			movies = append(movies, model.ParsedMovie{
				Title:    match[1],
				Year:     match[2],
				Language: match[3],
			})
		}
	}

	return movies
}
func GetMovieSuggestionsList(req model.MovieRequest, page int) ([]map[string]interface{}, error) {
	config := openai.DefaultConfig(os.Getenv("DEEPSEEK_API_KEY"))
	config.BaseURL = os.Getenv("DEEPSEEK_API_BASE_URL")
	client := openai.NewClientWithConfig(config)

	startIndex := (page - 1) * 10
	endIndex := startIndex + 10

	prompt := fmt.Sprintf(`
Given:
Hero: %s
Genre: %s
Language: %s
Reception: %s
Time of day: %s

Return EXACTLY 100 movies in this STRICT format:
"Title (Year) [ISO_Language_Code], Title (Year) [ISO_Language_Code], ..."

STRICT RULES:
1. RECEPTION CRITERIA:
   - "hit": >7/10 ratings AND box office success (top 25%% of year)
   - "flop": <5/10 ratings OR bottom 25%% box office
   - "underrated": >7/10 ratings BUT bottom 50%% box office
   - "overrated": <5/10 ratings DESPITE top 50%% box office
   - Custom terms: Interpret literally with clear justification

2. LANGUAGE REQUIREMENTS:
   - Primary language: %s (auto-convert to ISO code)
   - Dubbed versions ONLY if specified in request

3. CONTENT VALIDATION:
   - ONLY include released films (NO upcoming/announced projects)
   - NEVER include placeholder titles (e.g., "Actor 25", "NBK109")
   - Verify EVERY title exists on IMDb/TMDB

4. FORMATTING:
   - EXACT structure: "Title (Year) [code]"
   - ONLY comma separation (no line breaks/numbers)
   - ISO codes:
     • Telugu → te
     • Hindi → hi
     • Tamil → ta
     • Malayalam → ml
     • Kannada → kn
     • English → en
     • Others → First 2 letters (Japanese→ja)

5. SORTING:
   - Newest → oldest
   - Then by: 
     • "hit"/"overrated": Box office revenue (high→low)
     • "flop"/"underrated": Rating (high→low)

6. QUALITY CONTROL:
   - MUST reach 100 verified entries
   - Include obscure films if needed
   - NO commentary/explanations

7. DAY TIME SPECIFICITY:
	- If "time of day" is specified "Yes" consider the following:
		- "Morning" (6 AM - 12 PM): Motivational, feel-good films
		- "Afternoon" (12 PM - 6 PM): Comedy, family-friendly films
		- "Evening" (6 PM - 12 AM): Action, thriller, or drama
		- "Night" (12 AM - 6 AM): Horror, suspense, Romantic or mystery films

FAILURE CASES TO AVOID:
- "Ante Sundaraniki" for action
- "Devara (2024)" for underrated
- "NTR 31", "Prabhas 25" etc.
- Mixed reception films

Return ONLY movies %d to %d in the list.
`,
		req.Hero, req.Genre, req.Language, req.Reception, req.TimeOfDay,
		req.Language, startIndex+1, endIndex)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from DeepSeek")
	}

	// Parse DeepSeek response
	raw := resp.Choices[0].Message.Content
	parsedMovies := parseFormattedMovies(raw)
	fmt.Println("Parsed movies:", parsedMovies)

	uniqueMovies := make(map[string]bool)
	var results []map[string]interface{}

	for _, m := range parsedMovies {
		key := fmt.Sprintf("%s (%s)", m.Title, m.Year)
		if _, exists := uniqueMovies[key]; exists {
			continue // Skip duplicates
		}
		data, err := fetchMovieFromTMDB(m.Title, m.Language, m.Year)
		if err == nil && data != nil {
			uniqueMovies[key] = true
			results = append(results, data)
		}
	}

	return results, nil
}
