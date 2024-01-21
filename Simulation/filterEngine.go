package Simulation

import (
	"log"
)
import "github.com/TredingInGo/AutomationService/strategy"

type stockList struct {
	Symbol string `json:"symbol"`
}

var StockBook = make(map[string]string)

func PopulateStockTokens() {
	strategy.PopuletInstrumentsList()
	stockList := GetStockSymbolList()
	for i := 0; i < len(stockList); i++ {
		StockBook[stockList[i].Symbol] = strategy.GetToken(stockList[i].Symbol, "NSE")
	}
}

func GetStockSymbolList() []stockList {
	db := Connect()
	rows, err := db.Query(`SELECT * FROM "History"."StockList"`)
	if err != nil {
		log.Fatalf("Error querying data from the table: %v", err)
	}
	defer rows.Close()

	var stocks []stockList

	// Iterate through the rows and populate the array
	for rows.Next() {
		var stocksSymbol stockList
		err := rows.Scan(
			&stocksSymbol.Symbol,
		)
		if err != nil {
			log.Fatal(err)
		}
		stocks = append(stocks, stocksSymbol)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return stocks
}

//func FilterEngine() {
//	PopulateStockTokens()
//	var filteredStocks []string
//	var apiClient *smartapi.Client
//	for key, value := range StockBook {
//
//	}
//}
//
//func getHistoryData(client *smartapi.Client, symbolToken string, timeFrame string) {
//	tempTime := time.Now()
//	toDate := tempTime.Format("2006-01-02 15:04")
//	for i := 1; i <= 2; i++ {
//		fromDate := tempTime.Add(time.Hour * 24 * -5).Format("2006-01-02 15:04")
//		tempTime = tempTime.Add(time.Hour * 24 * -5)
//		tempHistoryData, _ := client.GetCandleData(smartapigo.CandleParams{
//			Exchange:    "NSE",
//			SymbolToken: symbolToken,
//			Interval:    timeFrame,
//			FromDate:    fromDate,
//			ToDate:      toDate,
//		})
//	}
//
//	return tempHistoryData
//}
