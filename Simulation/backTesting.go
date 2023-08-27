package Simulation

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math"
)

type tradeResult struct {
	amount       float64
	profitCount  int64
	lossCount    int64
	totalTrade   int64
	maxAmount    float64
	minAmount    float64
	stockSymbol  int64
	timeDuration int64
	smaLow       int
	smaHigh      int
	rsiPeriod    int
}
type tradeReport struct {
	amount      float64
	profit      float64
	loss        float64
	stratgyName string
	entry       float64
	tp          float64
	sl          float64
}
type tradeRecord struct {
	entryTimestamp string
	entryPrice     float64 // New field to store entry price
	exitTimestamp  string
	exitPrice      float64 // New field to store exit price
	profit         bool
}

type OHLC struct {
	Timestamp        string  `json:"timestamp"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Close            float64 `json:"close"`
	Volume           float64 `json:"volume"`
	ID               int64   `json:"id"`
	TimeFrameSeconds int64   `json:"timeFrameSeconds"`
}

var amount float64
var profitCount int64
var lossCount int64
var totalTrade int64
var maxAmount float64
var minAmount float64
var trades []tradeRecord
var trade []tradeReport

func DoBackTest(db *sql.DB) {

	ohlcData := GetData(db)
	//for i := 24; i < 100; i++ {
	//	for j := 20; j < i; j++ {
	//		for k := 32; k < 90; k++ {
	//			executetest(ohlcData, j, i, k, db)
	//
	//		}
	//	}
	//
	//}
	breakoutStretgy(ohlcData)
	saveTradeReport(trade, db)
}

func GetData(db *sql.DB) []OHLC {
	rows, err := db.Query(`SELECT * FROM "History"."OHLCData" WHERE id = 2885 AND timeframeinseconds = 900  order by timestamp limit 30000`)
	if err != nil {
		log.Fatalf("Error querying data from the table: %v", err)
	}
	defer rows.Close()

	var ohlcDataArray []OHLC

	// Iterate through the rows and populate the array
	for rows.Next() {
		var ohlc OHLC
		err := rows.Scan(
			&ohlc.ID,
			&ohlc.TimeFrameSeconds,
			&ohlc.Open,
			&ohlc.High,
			&ohlc.Low,
			&ohlc.Close,
			&ohlc.Timestamp,
			&ohlc.Volume,
		)
		if err != nil {
			log.Fatal(err)
		}
		ohlcDataArray = append(ohlcDataArray, ohlc)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ohlcDataArray

}

func executetest(ohlcData []OHLC, smaLow int, smaHigh int, rsiPeriod int, db *sql.DB) {
	amount = 100000.0
	maxAmount = amount
	minAmount = amount
	profitCount = 0
	lossCount = 0
	var tradeHistory tradeResult

	sma5 := CalculateSMA(ohlcData, smaLow)

	sma20 := CalculateSMA(ohlcData, smaHigh)

	k := 2.0
	ubb, lbb := CalculateBollingerBands(ohlcData, 20, k)

	for i := 100; i < len(ohlcData)-102; i++ {
		rsi := CalculateRSI(ohlcData, i, rsiPeriod)
		buySignal := false
		trend := "none"
		if ohlcData[i].Close > ubb[i-20] {
			trend = "up"
		} else if ohlcData[i].Close < lbb[i-20] {
			trend = "down"
		}
		if sma5[i] > sma20[i] && trend == "up" {

			if rsi > 60 && ohlcData[i].Close >= ohlcData[i-1].High {
				buySignal = true
			}
		}
		if buySignal {

			sl := ohlcData[i-1].Low
			bp := ohlcData[i].Close
			tp := (bp - sl) + bp
			PlaceBuyOrder(ohlcData, i, sl, tp, bp)
			fmt.Printf("Buy at %.2f, SL at %.2f, TP at %.2f\n", ohlcData[i].Close, sl, tp)
		}
	}
	tradeHistory.amount = amount
	tradeHistory.totalTrade = totalTrade
	tradeHistory.maxAmount = maxAmount
	tradeHistory.minAmount = minAmount
	tradeHistory.profitCount = profitCount
	tradeHistory.lossCount = lossCount
	tradeHistory.stockSymbol = ohlcData[0].ID
	tradeHistory.timeDuration = ohlcData[0].TimeFrameSeconds
	tradeHistory.smaLow = smaLow
	tradeHistory.smaHigh = smaHigh
	tradeHistory.rsiPeriod = rsiPeriod
	insertTradeSummery(tradeHistory, db)

}

func insertTradeSummery(trades tradeResult, db *sql.DB) {
	insertQuery := `INSERT INTO "History"."TradeSummary" (
                        amount, profitCount, lossCount, totalTrade,
						maxAmount, minAmount, stockSymbol, timeDuration,
						smaLow, smaHigh, rsiPeriod
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.Exec(insertQuery,
		trades.amount,
		trades.profitCount,
		trades.lossCount,
		trades.totalTrade,
		trades.maxAmount,
		trades.minAmount,
		trades.stockSymbol,
		trades.timeDuration,
		trades.smaLow,
		trades.smaHigh,
		trades.rsiPeriod,
	)

	if err != nil {
		fmt.Println("Error executing INSERT query:", err)
	}
}

// 2nd stretgy ...........

func breakoutStretgy(data []OHLC) {
	amount = 300000.0
	maxAmount = amount
	minAmount = amount
	profitCount = 0
	lossCount = 0

	high, low := calculateLowHigh(data, 4)

	for i := 4; i < len(data); i++ {
		// buying logic
		if amount <= 0 {
			return
		}
		if i%24 == 0 && i+4 < len(data) {
			high, low = calculateLowHigh(data, i+4)
		}
		weag := calculateWeightedPercentage(data[i])

		if data[i].Close > high && data[i].Close < data[i].Open && weag < 17 {
			triggerBuyOrder(data, data[i], i, high)
		} else if data[i].Close < high && data[i].Close > data[i].Open && weag < 17 {
			triggerSellOrder(data, data[i], i, low)
		}
	}
}

func calculateLowHigh(data []OHLC, index int) (high float64, low float64) {
	high = 0.0
	low = 100000.0
	for it := index - 1; it >= index-4; it-- {
		high = math.Max(high, data[it].High)
		low = math.Min(low, data[it].Low)
	}
	return high, low
}

func calculateWeightedPercentage(candle OHLC) float64 {
	weightedHigh := math.Abs(candle.High - candle.Low)
	weightedClose := math.Abs(candle.Close - candle.Open)
	weightedPercentage := (weightedClose * 100) / weightedHigh
	return weightedPercentage
}
func triggerBuyOrder(data []OHLC, candle OHLC, index int, high float64) {
	for i := index; i < len(data); i++ {
		if data[i].Low < high {
			return
		}
		if data[i].High >= candle.High {
			PlaceBuyOrder(data, i, candle.Low, data[i].High+2*(data[i].High-candle.Low), data[i].High)
		}
	}
}

func triggerSellOrder(data []OHLC, candle OHLC, index int, low float64) {
	for i := index; i < len(data); i++ {
		if data[i].High > low {
			return
		}
		if data[i].Low <= candle.Low {
			//PlaceSellOrder(data, i, candle.High, data[i].Low-2*(candle.High-data[i].Low), data[i].Low)
		}
	}
}

func saveTradeReport(trades []tradeReport, db *sql.DB) {
	for i := 0; i < len(trades); i++ {
		insertSQL := `
		INSERT INTO "History"."TradeReport" (amount, profit, loss, strategyName, entry, tp, sl)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`
		_, err := db.Exec(insertSQL,
			trades[i].amount,
			trades[i].profit,
			trades[i].loss,
			trades[i].stratgyName,
			trades[i].entry,
			trades[i].tp,
			trades[i].sl,
		)

		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("Inserted row with ID: %d\n", insertedID)
	}
}
