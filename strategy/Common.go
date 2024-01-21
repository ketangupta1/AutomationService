package strategy

import (
	"database/sql"
	smartapigo "github.com/TredingInGo/smartapi"
	"log"
	"time"
)

func LoadStockList(db *sql.DB) []Symbols {
	rows, err := db.Query(`SELECT * FROM "History"."Intraday"`)
	if err != nil {
		log.Fatalf("Error querying data from the table: %v", err)
	}
	defer rows.Close()

	var stockList []Symbols

	// Iterate through the rows and populate the array
	for rows.Next() {
		var stock Symbols
		err := rows.Scan(
			&stock.Symbol,
			&stock.Token,
		)
		if err != nil {
			log.Fatal(err)
		}
		stockList = append(stockList, stock)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return stockList

}

func GetStockTick(client *smartapigo.Client, symbolToken string, timeFrame string) []smartapigo.CandleResponse {
	tempTime := time.Now()
	toDate := tempTime.Format("2006-01-02 15:04")
	fromDate := tempTime.Add(time.Hour * 24 * -5).Format("2006-01-02 15:04")
	tempTime = tempTime.Add(time.Hour * 24 * -5)
	tempHistoryData, _ := client.GetCandleData(smartapigo.CandleParams{
		Exchange:    "NSE",
		SymbolToken: symbolToken,
		Interval:    timeFrame,
		FromDate:    fromDate,
		ToDate:      toDate,
	})
	return tempHistoryData
}
