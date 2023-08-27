package Simulation

import (
	"database/sql"
	"fmt"
	smartapigo "github.com/TredingInGo/smartapi"
	"strconv"
	"time"
)

func PrepareData(db *sql.DB, client *smartapigo.Client) {
	symbolToken := 2885
	tempTime := time.Now()
	for j := 0; j < 730; j++ {
		toDate := tempTime.Format("2006-01-02 15:04")
		fromDate := tempTime.Add(time.Hour * 24 * -5).Format("2006-01-02 15:04")
		tempTime = tempTime.Add(time.Hour * 24 * -5)
		tempHistoryData, _ := client.GetCandleData(smartapigo.CandleParams{
			Exchange:    "NSE",
			SymbolToken: strconv.Itoa(symbolToken),
			Interval:    "FIFTEEN_MINUTE",
			FromDate:    fromDate,
			ToDate:      toDate,
		})

		for _, tempHistoryData := range tempHistoryData {
			// Prepare the INSERT statement
			insertQuery := `
                INSERT INTO "History"."OHLCData" (id, timeframeinseconds, open, high, low, close, timestamp, volume) 
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

			// Execute the INSERT statement
			_, err := db.Exec(insertQuery,
				symbolToken,
				900,
				tempHistoryData.Open,
				tempHistoryData.High,
				tempHistoryData.Low,
				tempHistoryData.Close,
				tempHistoryData.Timestamp,
				tempHistoryData.Volume,
			)

			if err != nil {
				fmt.Println("Error executing INSERT query:", err)
				return
			}
		}

	}
}
