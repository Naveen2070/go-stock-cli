package csvloader

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/Naveen2070/go-stock-cli/models"
)

func Load(filename string) ([]models.Stock, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	rows = slices.Delete(rows, 0, 1)

	var stocks []models.Stock
	for _, row := range rows {
		ticker := row[0]
		gap, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}
		openingPrice, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}
		stocks = append(stocks, models.Stock{Ticker: ticker, Gap: gap, OpeningPrice: openingPrice})
	}

	return stocks, nil
}
