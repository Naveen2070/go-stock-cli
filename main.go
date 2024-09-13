package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

type Stock struct {
	ticker       string
	gap          float64
	openingPrice float64
}

func load(filename string) ([]Stock, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rows = slices.Delete(rows, 0, 1)

	var stocks []Stock

	for _, row := range rows {
		ticker := row[0]
		gap, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		openingPrice, err := strconv.ParseFloat(row[2], 64)

		if err != nil {
			fmt.Println(err)
			continue
		}

		stocks = append(stocks, Stock{ticker, gap, openingPrice})
	}

	return stocks, nil
}

var accountBalance = 100.00
var lossTolerance = 0.02
var maxLossPerTrade = accountBalance * lossTolerance
var profitPercentage = 0.08

type Position struct {
	EntryPrice       float64
	Shares           int
	TakeProfilePrice float64
	StopLossPrice    float64
	Profit           float64
}

func caculate(gapPercent, openingPrice float64) Position {
	closingPrice := openingPrice / (1 + gapPercent)
	gapValue := closingPrice - openingPrice
	profitFromGap := gapValue * profitPercentage

	stopLoss := openingPrice - profitFromGap
	takeProfit := closingPrice + profitFromGap

	shares := int(maxLossPerTrade / math.Abs(stopLoss-openingPrice))

	profit := math.Abs(openingPrice-takeProfit) * float64(shares)
	profit = math.Round(profit*100) / 100

	return Position{
		EntryPrice:       math.Round(openingPrice*100) / 100,
		Shares:           shares,
		TakeProfilePrice: math.Round(takeProfit*100) / 100,
		StopLossPrice:    math.Round(stopLoss*100) / 100,
		Profit:           math.Round(profit*100) / 100,
	}
}

type Selection struct {
	Ticker string
	Position
	Articles []Article
}

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

type Article struct {
	PublishedOn time.Time
	Headline    string
}

const (
	url       = "https://seeking-alpha.p.rapidapi.com/news/v2/list-by-symbol?size=5&id="
	apiHeader = "X-RapidAPI-Key"
	apiKey    = "61550de255msh906da11198a49e6p1fcb7djsnad675fd4be5d"
)

func fetchNews(ticker string) ([]Article, error) {
	req, err := http.NewRequest(http.MethodGet, url+ticker, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add(apiHeader, apiKey)
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errorMsg := fmt.Errorf("Request failed with status code: %d", resp.StatusCode)
		return nil, errorMsg
	}

	if resp.Body == nil {
		fmt.Println("Response body is nil")
		return nil, err
	}

	res := &seekingAlphaResponse{}
	json.NewDecoder(resp.Body).Decode(res)

	var articles []Article

	for _, news := range res.Data {
		art := Article{
			PublishedOn: news.Attributes.PublishedOn,
			Headline:    news.Attributes.Title,
		}

		articles = append(articles, art)
	}

	defer resp.Body.Close()
	return articles, nil
}

func Deliver(filePath string, selections []Selection) error {

	f, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(selections)

	if err != nil {
		return fmt.Errorf("Error encoding selections: %v", err)
	}

	return nil
}

func main() {
	stocks, err := load("./opg.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	stocks = slices.DeleteFunc(stocks, func(stock Stock) bool {
		return math.Abs(stock.gap) < 0.1
	})

	var selections []Selection

	// var waitGroup sync.WaitGroup
	selectionChannel := make(chan Selection, len(stocks))

	for _, stock := range stocks {
		// waitGroup.Add(1)

		go func(stock Stock, selectionChannel chan<- Selection) {
			// defer waitGroup.Done()
			position := caculate(stock.gap, stock.openingPrice)

			articles, err := fetchNews(stock.ticker)

			if err != nil {
				log.Printf("Error fetching news for %s: %v", stock.ticker, err)
				return
			} else {
				log.Printf("Successfully fetched %d news for %s", len(articles), stock.ticker)
			}

			selected := Selection{
				Ticker:   stock.ticker,
				Position: position,
				Articles: articles,
			}

			selectionChannel <- selected
		}(stock, selectionChannel)

	}

	// waitGroup.Wait()

	for sel := range selectionChannel {
		selections = append(selections, sel)

		if len(selections) == len(stocks) {
			close(selectionChannel)
		}
	}

	outputPath := "./Analysis-Result.json"
	err = Deliver(outputPath, selections)

	if err != nil {
		log.Printf("Error delivering selections: %v", err)
		return
	}

	log.Printf("Successfully delivered selections to %s \n", outputPath)
}
