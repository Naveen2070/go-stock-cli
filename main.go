package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
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
}

func main() {
	stocks, err := load("./opg.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	slices.DeleteFunc(stocks, func(stock Stock) bool {
		return math.Abs(stock.gap) < 0.1
	})

	var selections []Selection

	for _, stock := range stocks {
		position := caculate(stock.gap, stock.openingPrice)

		selected := Selection{
			Ticker:   stock.ticker,
			Position: position,
		}

		selections = append(selections, selected)
	}
}
