package Simulation

import (
	"database/sql"
	"encoding/json"
	"fmt"
	smartapigo "github.com/TredingInGo/smartapi"
	"strconv"
	"time"
)

func PrepareData(db *sql.DB, client *smartapigo.Client) {
	symbolToken := 2885
	tempTime := time.Now()
	historyData := make([]smartapigo.CandleResponse, 0)
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
		historyData = append(historyData, tempHistoryData...)
	}
	insertQuery := `
        INSERT INTO "History"."HistoryData" (id, timeframeinseconds, ohlc) 
        VALUES ($1, $2, $3)`
	data := struct {
		Ohlc []smartapigo.CandleResponse `json:"ohlc""`
	}{Ohlc: historyData}

	ohlcData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling OHLC data:", err)
		return
	}

	_, err = db.Exec(insertQuery, symbolToken, 900, ohlcData)
	if err != nil {
		fmt.Println("Error executing INSERT query:", err)
		return
	}

}
