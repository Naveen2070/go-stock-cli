package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Naveen2070/go-stock-cli/models"
)

const (
	apiURL    = "https://seeking-alpha.p.rapidapi.com/news/v2/list-by-symbol?size=5&id="
	apiHeader = "X-RapidAPI-Key"
	apiKey    = "61550de255msh906da11198a49e6p1fcb7djsnad675fd4be5d"
)

type attribute struct {
	PublishedOn time.Time `json:"publishon"`
	Title       string    `json:"title"`
}

type seekingAlphaNews struct {
	Attributes attribute `json:"attributes"`
}

type seekingAlphaResponse struct {
	Data []seekingAlphaNews `json:"data"`
}

func FetchNews(ticker string) ([]models.Article, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL+ticker, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add(apiHeader, apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var res seekingAlphaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var articles []models.Article
	for _, news := range res.Data {
		articles = append(articles, models.Article{
			PublishedOn: news.Attributes.PublishedOn.String(),
			Headline:    news.Attributes.Title,
		})
	}
	return articles, nil
}
