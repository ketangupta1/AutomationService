package Simulation

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math"
)

type OHLC struct {
	Timestamp string  `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
}

var amount float64
var profitCount int64
var lossCount int64
var totalTrade int64
var maxAmount float64
var minAmount float64

func DoBackTest(db *sql.DB) {
	ohlcData := getData(db)
	amount = 100000.0
	maxAmount = amount
	minAmount = amount
	profitCount = 0
	lossCount = 0

	sma5 := calculateSMA(ohlcData, 5)

	ema20 := calculateEMA(ohlcData, 20)

	k := 2.0
	ubb, lbb := calculateBollingerBands(ohlcData, 20, k)

	for i := 20; i < len(ohlcData)-20; i++ {
		rsi := calculateRSI(ohlcData, i)
		buySignal := false
		trend := "none"
		if ohlcData[i].Close > ubb[i-20] {
			trend = "up"
		} else if ohlcData[i].Close < lbb[i-20] {
			trend = "down"
		}
		if sma5[i] > ema20[i] && trend == "up" {

			if rsi > 60 && ohlcData[i].Close >= ohlcData[i-1].High {
				buySignal = true
			}
		}
		if buySignal {

			sl := ohlcData[i-1].Low
			bp := ohlcData[i-1].High
			tp := 2*(bp-sl) + bp
			placeOrder(ohlcData, i, sl, tp, bp)
			fmt.Printf("Buy at %.2f, SL at %.2f, TP at %.2f\n", ohlcData[i].Close, sl, tp)
		}
	}

	fmt.Printf("Amount: %.2f\n", amount)
	fmt.Printf("Profit Count: %d\n", profitCount)
	fmt.Printf("Loss Count: %d\n", lossCount)
	fmt.Printf("Total Trades: %d\n", totalTrade)
	fmt.Printf("Maximum Amount: %.2f\n", maxAmount)
	fmt.Printf("Minimum Amount: %.2f\n", minAmount)

}

func calculateSMA(data []OHLC, period int) []float64 {
	sma := make([]float64, len(data)-period+1)

	for i := 0; i <= len(data)-period; i++ {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += data[j].Close
		}
		sma[i] = sum / float64(period)
	}

	return sma
}

func calculateEMA(data []OHLC, period int) []float64 {
	ema := make([]float64, len(data)-period+1)

	sma := calculateSMA(data, period)

	for i := 0; i < period; i++ {
		ema[i] = sma[i]
	}

	smoothingFactor := 2.0 / (float64(period) + 1)
	for i := period; i < len(data); i++ {
		ema[i-period+1] = (data[i].Close-sma[i-period])*smoothingFactor + ema[i-period]
	}

	return ema
}

func calculateRSI(data []OHLC, index int) float64 {
	rsiPeriod := 14
	if index < rsiPeriod {
		return 0.0 // RSI not available for the first few data points
	}

	gain := 0.0
	loss := 0.0

	// Calculate average gain and loss
	for i := index; i > index-rsiPeriod; i-- {
		diff := data[i].Close - data[i-1].Close
		if diff >= 0 {
			gain += diff
		} else {
			loss -= diff
		}
	}
	avgGain := gain / float64(rsiPeriod)
	avgLoss := loss / float64(rsiPeriod)

	// Calculate RSI
	if avgLoss == 0 {
		return 100.0 // Avoid division by zero
	} else {
		rs := avgGain / avgLoss
		return 100.0 - (100.0 / (1.0 + rs))
	}
}

func getData(db *sql.DB) []OHLC {
	rows, err := db.Query(`SELECT ohlc FROM "History"."HistoryData" WHERE id = 2885 AND timeframeinseconds = 900`)
	if err != nil {
		log.Fatalf("Error querying data from the table: %v", err)
	}
	defer rows.Close()

	var data struct {
		DATA []OHLC `json:"ohlc"`
	}
	// Iterate through the rows and parse the JSON data
	for rows.Next() {
		var ohlcJSON string

		if err := rows.Scan(&ohlcJSON); err != nil {
			log.Fatalf("Error scanning data from the row: %v", err)
		}

		if err := json.Unmarshal([]byte(ohlcJSON), &data); err != nil {
			log.Fatalf("Error unmarshalling JSON data: %v", err)
		}

		// Append the OHLC data to the slice

	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating through the rows: %v", err)
	}

	return data.DATA

}

func calculateBollingerBands(data []OHLC, period int, k float64) ([]float64, []float64) {
	sma := calculateSMA(data, period)
	sd := calculateSD(data, period)
	ubb := make([]float64, len(data)-period+1)
	lbb := make([]float64, len(data)-period+1)

	for i := period; i < len(data); i++ {
		ubb[i-period] = sma[i-period] + (sd[i-period] * k)
		lbb[i-period] = sma[i-period] - (sd[i-period] * k)
	}

	return ubb, lbb
}

func calculateSD(data []OHLC, period int) []float64 {
	sd := make([]float64, len(data)-period)

	for i := period; i < len(data)-period; i++ {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += (data[j].Close - data[i-1].Close) * (data[j].Close - data[i-1].Close)
		}
		mean := sum / float64(period)
		sd[i-period] = math.Sqrt(mean)
	}
	return sd
}

func placeOrder(ohlcData []OHLC, index int, sl, takeProfit, buyPrice float64) {

	// Loop from the current index + 1 to the end of the OHLC data
	totalTrade++
	positionSize := calculatePositionSize(buyPrice, sl)
	fmt.Printf("PostionSize: %v", positionSize)
	for i := index + 1; i < len(ohlcData); i++ {
		if ohlcData[i].Low <= sl {
			// Stop Loss triggered, calculate loss
			loss := (sl - buyPrice) * positionSize
			lossCount++
			amount += loss
			fmt.Printf("Trade result: Loss %.2f\n", loss)
			minAmount = math.Min(minAmount, amount)
			return
		} else if ohlcData[i].High >= takeProfit {
			// Take Profit triggered, calculate profit
			profit := (takeProfit - buyPrice) * positionSize
			profitCount++
			amount += profit
			maxAmount = math.Max(maxAmount, amount)
			fmt.Printf("Trade result: Profit %.2f\n", profit)
			return
		}
	}

	// If the loop completes without hitting Stop Loss or Take Profit
	fmt.Println("Trade result: No Stop Loss or Take Profit triggered.")

	// If the loop completes without hitting Stop Loss or Take Profit
	return
}

func calculatePositionSize(buyPrice, sl float64) float64 {
	maxRiskPercent := 0.02 // 2% maximum risk allowed
	maxRiskAmount := amount * maxRiskPercent
	riskPerShare := math.Max(1, buyPrice-sl)
	positionSize := maxRiskAmount / riskPerShare
	for positionSize*maxRiskAmount > amount {
		positionSize--
	}

	return positionSize
}
