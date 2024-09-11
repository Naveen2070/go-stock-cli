package main

import (
	"encoding/csv"
	"fmt"
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
func main() {
	stocks, err := load("./opg.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stocks)
}
