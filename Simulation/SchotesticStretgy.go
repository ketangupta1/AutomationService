package Simulation

import (
	"database/sql"
	"log"
	"math"
)

type tradeReportSch struct {
	amountToTest float64
	profit       float64
	loss         float64
	stratgyName  string
	kVal         float64
	period       int
	tradeCount   int
	maxDrawDown  float64
	maxProfit    float64
}

var amountToTest float64
var tradeCount int
var maxDrawDown float64
var maxProfit float64
var tradesHistory []tradeReportSch

func SchRun(db *sql.DB) {
	ohlc := GetData(db)
	for i := 10; i <= 25; i++ {
		execute(ohlc, i)
	}
	saveTradesSch(db)

}

func execute(ohlc []OHLC, period int) {
	for sValue := 75.0; sValue <= 85; sValue++ {
		amountToTest = 300000
		tradeCount = 0
		maxDrawDown = amountToTest
		maxProfit = 0
		for i := period + 1; i < len(ohlc); i++ {
			data := CalculateStochastic(ohlc, i, period)
			if data.K > sValue {
				buyOrder(ohlc, i, sValue, period)
			}
			if data.K < 20 {
				//sellOrder(ohlc, i, sValue, period)
			}
		}

	}
}

func buyOrder(data []OHLC, index int, sValue float64, period int) {
	buyPrice := data[index].Close
	sl := data[index].Low
	tp := buyPrice + ((buyPrice - sl) * 2)
	tradeCount++
	postionSize := calculatePositionSizeSch(buyPrice, sl)
	for i := index + 1; i < len(data); i++ {
		if data[i].Low <= sl {
			profit := (buyPrice - sl) * postionSize
			amountToTest -= profit
			maxDrawDown = math.Min(maxDrawDown, amountToTest)
			tradesHistory = append(tradesHistory,
				tradeReportSch{amountToTest, 0, profit, "scho", sValue, period, tradeCount, maxDrawDown, maxProfit})

			return
		} else if data[i].High >= tp {
			profit := (tp - buyPrice) * postionSize
			amountToTest += profit
			maxProfit = math.Max(maxProfit, amountToTest)
			tradesHistory = append(tradesHistory,
				tradeReportSch{amountToTest, profit, 0, "scho", sValue, period, tradeCount, maxDrawDown, maxProfit})

			return
		}
	}
}

func sellOrder(data []OHLC, index int, sVal float64, period int) {
	buyPrice := data[index].Close
	sl := data[index].High
	tp := buyPrice - ((sl - buyPrice) * 2)
	tradeCount++
	postionSize := calculatePositionSizeSch(buyPrice, sl)
	for i := index + 1; i < len(data); i++ {
		if data[i].High >= sl {
			amountToTest -= (sl - buyPrice) * postionSize
			maxDrawDown = math.Min(maxDrawDown, amountToTest)
			return
		} else if data[i].High >= tp {
			amountToTest += (buyPrice - tp) * postionSize
			maxProfit = math.Max(maxProfit, amountToTest)
			return
		}
	}

}

func saveTradesSch(db *sql.DB) {
	for i := 0; i < len(tradesHistory); i++ {
		insertSQL := `
		INSERT INTO "History"."TradeReportSch" (amounttotest, profit, loss, strategyname, kval, periods, tradecount, maxdrawdown, maxprofit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;`
		_, err := db.Exec(insertSQL,
			tradesHistory[i].amountToTest,
			tradesHistory[i].profit,
			tradesHistory[i].loss,
			tradesHistory[i].stratgyName,
			tradesHistory[i].kVal,
			tradesHistory[i].period,
			tradesHistory[i].tradeCount,
			tradesHistory[i].maxDrawDown,
			tradesHistory[i].maxProfit,
		)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func calculatePositionSizeSch(buyPrice, sl float64) float64 {
	if amountToTest <= 0 {
		return 0
	}
	maxRiskPercent := 0.05 // 2% maximum risk allowed
	maxRiskAmount := amountToTest * maxRiskPercent
	riskPerShare := math.Max(1, math.Abs(buyPrice-sl))
	positionSize := maxRiskAmount / riskPerShare
	for positionSize*maxRiskAmount > amountToTest {
		positionSize--
	}

	return positionSize
}
