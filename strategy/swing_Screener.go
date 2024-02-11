package strategy

import (
	"database/sql"
	"fmt"
	smartapigo "github.com/TredingInGo/smartapi"
	"math"
	"strconv"
)

func SwingScreener(client *smartapigo.Client, db *sql.DB) {
	stockList := LoadStockList(db)
	fmt.Printf("*************** LIST FOR SWING TRADING *****************")

	for _, stock := range stockList {
		ExecuteScreener(stock.Token, stock.Symbol, client)
		return
	}
}

func ExecuteScreener(symbol, stockToken string, client *smartapigo.Client) {
	fmt.Printf("\n%v", symbol)
	data := GetStockTickForSwing(client, stockToken, "FOUR_HOUR")
	for i := len(data) - 1; i >= len(data)-6; i-- {
		fmt.Printf("\n%v\n", data[i])
	}
	if len(data) == 0 {
		return
	}
	PopulateIndicators(data, stockToken)
	order := TrendFollowingRsiForSwing(data, stockToken, symbol)
	if order.OrderType == "None" {
		return
	}
	orderParams := SetOrderParamsForSwing(order, stockToken, symbol)
	countStock := 1
	fmt.Printf("\n                   STOCK No: %v                        \n", countStock)
	fmt.Printf("\n=========================================================\n")
	fmt.Printf("\nSTOCK NAME -  %v\n", symbol)
	fmt.Printf("SPOT PRICE - %v\n", order.Spot)
	fmt.Printf("STOP LOSS -  %v\n", order.Sl)
	fmt.Printf("TARGET -      %v\n\n", order.Tp)
	fmt.Printf("Order Params -      %v\n\n", orderParams)
	fmt.Printf("\n=========================================================\n\n")
	countStock++
}

func TrendFollowingRsiForSwing(data []smartapigo.CandleResponse, token, symbol string) ORDER {
	idx := len(data) - 1
	sma5 := sma[token+"5"][idx]
	sma8 := sma[token+"8"][idx]
	adx14 := adx[token]
	rsi := rsi[token]
	var order ORDER
	order.OrderType = "None"
	swingLow := GetSwingLow(data, 20)
	if adx14.Adx[idx] >= 25 && adx14.PlusDi[idx] > adx14.MinusDi[idx] && sma5 > sma8 && sma8 > sma[token+"13"][idx] && sma[token+"13"][idx] > sma[token+"21"][idx] && rsi[idx] < 70 && rsi[idx] > 55 && rsi[idx-2] < rsi[idx] {
		order = ORDER{
			Spot:      data[idx].High + 0.05,
			Sl:        int(swingLow),
			Tp:        int(data[idx].High + 2*(data[idx].High-swingLow)),
			Quantity:  1,
			OrderType: "BUY",
		}
	}

	return order
}

func SetOrderParamsForSwing(order ORDER, token, symbol string) smartapigo.OrderParams {

	orderParams := smartapigo.OrderParams{
		Variety:          "GTT",
		TradingSymbol:    symbol + "-EQ",
		SymbolToken:      token,
		TransactionType:  order.OrderType,
		Exchange:         "NSE",
		OrderType:        "LIMIT",
		ProductType:      "GTT",
		Duration:         "DELIVERY",
		Price:            strconv.FormatFloat(order.Spot, 'f', 2, 64),
		SquareOff:        strconv.Itoa(order.Tp),
		StopLoss:         strconv.Itoa(order.Sl),
		Quantity:         strconv.Itoa(order.Quantity),
		TrailingStopLoss: strconv.Itoa(1),
	}
	return orderParams
}

func GetSwingLow(data []smartapigo.CandleResponse, day int) float64 {
	length := len(data)
	low := 1000000.0
	for i := length - 1; i > length-day-1; i-- {
		low = math.Min(low, data[i].Low)
	}
	return low
}
