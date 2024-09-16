package main

import (
	"log"
	"math"
	"slices"
	"sync"

	"github.com/Naveen2070/go-stock-cli/internal/calculations"
	"github.com/Naveen2070/go-stock-cli/internal/csvloader"
	"github.com/Naveen2070/go-stock-cli/internal/fetcher"
	"github.com/Naveen2070/go-stock-cli/internal/output"
	"github.com/Naveen2070/go-stock-cli/models"
)

func main() {
	stocks, err := csvloader.Load("./opg.csv")
	if err != nil {
		log.Fatal(err)
	}

	stocks = slices.DeleteFunc(stocks, func(stock models.Stock) bool {
		return math.Abs(stock.Gap) < 0.1
	})

	var selections []models.Selection
	var waitGroup sync.WaitGroup
	selectionChannel := make(chan models.Selection, len(stocks))

	for _, stock := range stocks {
		waitGroup.Add(1)
		go func(stock models.Stock, selectionChannel chan<- models.Selection) {
			defer waitGroup.Done()
			position := calculations.Calculate(stock.Gap, stock.OpeningPrice)

			articles, err := fetcher.FetchNews(stock.Ticker)
			if err != nil {
				log.Printf("Error fetching news for %s: %v", stock.Ticker, err)
				return
			} else {
				log.Printf("Successfully fetched %d news for %s", len(articles), stock.Ticker)
			}

			selected := models.Selection{
				Ticker:   stock.Ticker,
				Position: position,
				Articles: articles,
			}
			selectionChannel <- selected
		}(stock, selectionChannel)
	}

	go func() {
		waitGroup.Wait()
		close(selectionChannel)
	}()

	for sel := range selectionChannel {
		selections = append(selections, sel)
	}

	err = output.Deliver("./Analysis-Result.json", selections)
	if err != nil {
		log.Fatalf("Error delivering selections: %v", err)
	}

	log.Printf("Successfully delivered selections to %s \n", "./Analysis-Result.json")
}
