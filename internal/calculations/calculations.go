package calculations

import (
	"math"

	"github.com/Naveen2070/go-stock-cli/models"
)

var accountBalance = 100.00
var lossTolerance = 0.02
var maxLossPerTrade = accountBalance * lossTolerance
var profitPercentage = 0.08

func Calculate(gapPercent, openingPrice float64) models.Position {
	closingPrice := openingPrice / (1 + gapPercent)
	gapValue := closingPrice - openingPrice
	profitFromGap := gapValue * profitPercentage

	stopLoss := openingPrice - profitFromGap
	takeProfit := closingPrice + profitFromGap

	shares := int(maxLossPerTrade / math.Abs(stopLoss-openingPrice))
	profit := math.Abs(openingPrice-takeProfit) * float64(shares)
	profit = math.Round(profit*100) / 100

	return models.Position{
		EntryPrice:       math.Round(openingPrice*100) / 100,
		Shares:           shares,
		TakeProfilePrice: math.Round(takeProfit*100) / 100,
		StopLossPrice:    math.Round(stopLoss*100) / 100,
		Profit:           math.Round(profit*100) / 100,
	}
}
