package strategy

import (
	"database/sql"
	"fmt"
	smartapigo "github.com/TredingInGo/smartapi"
	"strconv"
	"sync"
	"time"
)

type OrderDetails struct {
	Spot      float64
	Tp        float64
	Sl        float64
	Quantity  int
	OrderType string
}

type Symbols struct {
	Symbol string `json:"symbol"`
	Token  string `json:"token"`
}

var flag = false
var amount = 0.0
var orderId = make(map[string]bool)
var myMutex sync.Mutex

func CloseSession(client *smartapigo.Client) {
	currentTime := time.Now()
	compareTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 0, 0, 0, currentTime.Location())

	if currentTime.After(compareTime) {
		client.Logout()
		fmt.Printf("Session closed ")
	}

}
func TrendFollowingStretgy(client *smartapigo.Client, db *sql.DB) {
	go CloseSession(client)
	stockList := LoadStockList(db)
	UpdateAmount(client)
	TrackOrders(client, "DUMMY")
	for {

		for _, stock := range stockList {

			if !flag {
				Execute(stock.Token, stock.Symbol, client)
			}

		}
		//time.Sleep(5 * time.Second)

	}
}

func Execute(symbol, stockToken string, client *smartapigo.Client) {
	data := GetStockTick(client, stockToken, "FIVE_MINUTE")
	if len(data) == 0 {
		return
	}
	lptParams := smartapigo.LTPParams{
		"NSE",
		symbol,
		stockToken,
	}
	ltp, _ := client.GetLTP(lptParams)
	dataToAppend := smartapigo.CandleResponse{
		Timestamp: time.Now(),
		Open:      ltp.Open,
		High:      ltp.High,
		Low:       ltp.Low,
		Close:     ltp.Close,
		Volume:    0,
	}
	data = append(data, dataToAppend)
	myMutex.Lock()
	PopulateIndicators(data, stockToken)

	myMutex.Unlock()
	UpdateAmount(client)
	order := TrendFollowingRsi(data, stockToken)
	if order.OrderType == "None" {
		return
	}
	if order.Quantity < 1 {

		return
	}
	orderParams := SetOrderParams(order, stockToken, symbol)
	fmt.Printf("\norder params:\n%v\n", orderParams)
	var orderRes smartapigo.OrderResponse
	myMutex.Lock()
	if flag == false {
		//flag = true
		orderRes, _ = client.PlaceOrder(orderParams)
		fmt.Printf("order response %v", orderRes)

	}

	myMutex.Unlock()
	UpdateAmount(client)
	TrackOrders(client, symbol)
	UpdateAmount(client)
}

func TrendFollowingRsi(data []smartapigo.CandleResponse, token string) ORDER {
	idx := len(data) - 1
	sma1 := 7
	sma2 := 22
	sma3 := sma[token+"3"][idx]
	sma5 := sma[token+"5"][idx]
	sma8 := sma[token+"8"][idx]
	adx14 := adx[token]
	rsi := rsi[token]
	var order ORDER
	order.OrderType = "None"
	if adx14.Adx[idx] >= 25 && adx14.PlusDi[idx] > adx14.MinusDi[idx] && sma3 > sma5 && sma5 > sma8 && sma8 > sma[token+"13"][idx] && sma[token+"13"][idx] > sma[token+"21"][idx] && rsi[idx] < 75 && rsi[idx] > 60 && rsi[idx-2] < rsi[idx] {
		fmt.Printf("order placed: trend following adx = %v \n", adx14.Adx[idx])
		order = ORDER{
			Spot:      data[idx].High + 0.05,
			Sl:        int(data[idx].High * 0.01),
			Tp:        int(data[idx].High * 0.02),
			Quantity:  CalculatePosition(data[idx].High, data[idx].High-data[idx].High*0.01),
			OrderType: "BUY",
		}
	} else if adx14.Adx[idx] >= 25 && adx14.PlusDi[idx] < adx14.MinusDi[idx] && sma3 < sma5 && sma5 < sma8 && sma8 < sma[token+"13"][idx] && sma[token+"13"][idx] < sma[token+"21"][idx] && rsi[idx] < 40 && rsi[idx] > 30 && rsi[idx-2] > rsi[idx] {
		fmt.Printf("order placed: trend following %v\n", adx14.Adx[idx])
		order = ORDER{
			Spot:      data[idx].Low - 0.05,
			Sl:        int(data[idx].High * 0.01),
			Tp:        int(data[idx].High * 0.02),
			Quantity:  CalculatePosition(data[idx].High, data[idx].High-data[idx].High*0.01),
			OrderType: "SELL",
		}
	} else if adx14.Adx[idx] >= 25 && adx14.PlusDi[idx] > adx14.MinusDi[idx] && sma[token+"7"][idx-1] < sma[token+"22"][idx-1] && sma[token+strconv.Itoa(sma1)][idx] > sma[token+strconv.Itoa(sma2)][idx] && rsi[idx] < 30 && rsi[idx-1] < rsi[idx] {
		fmt.Printf("order placed: trend reversal\n")
		order = ORDER{
			Spot:      data[idx].High + 0.05,
			Sl:        int(data[idx].High * 0.01),
			Tp:        int(data[idx].High * 0.02),
			Quantity:  CalculatePosition(data[idx].High, data[idx].High-data[idx].High*0.01),
			OrderType: "BUY",
		}
	} else if adx14.Adx[idx] >= 25 && adx14.PlusDi[idx] < adx14.MinusDi[idx] && sma[token+"7"][idx-1] > sma[token+"22"][idx-1] && sma[token+strconv.Itoa(sma1)][idx] < sma[token+strconv.Itoa(sma2)][idx] && rsi[idx] > 75 && rsi[idx-1] > rsi[idx] {
		fmt.Printf("order placed: trend reversal\n")
		order = ORDER{
			Spot:      data[idx].Low - 0.5,
			Sl:        int(data[idx].High * 0.01),
			Tp:        int(data[idx].High * 0.02),
			Quantity:  CalculatePosition(data[idx].Low, data[idx].High-data[idx].High*0.01),
			OrderType: "SELL",
		}
	}
	fmt.Printf("order placed: %v\n", order)
	return order
}

func SetOrderParams(order ORDER, token, symbol string) smartapigo.OrderParams {

	orderParams := smartapigo.OrderParams{
		Variety:          "ROBO",
		TradingSymbol:    symbol + "-EQ",
		SymbolToken:      token,
		TransactionType:  order.OrderType,
		Exchange:         "NSE",
		OrderType:        "LIMIT",
		ProductType:      "BO",
		Duration:         "DAY",
		Price:            strconv.FormatFloat(order.Spot, 'f', 2, 64),
		SquareOff:        strconv.Itoa(order.Tp),
		StopLoss:         strconv.Itoa(order.Sl),
		Quantity:         strconv.Itoa(order.Quantity),
		TrailingStopLoss: strconv.Itoa(order.Sl),
	}
	return orderParams
}
func UpdateAmount(client *smartapigo.Client) {
	RMS, _ := client.GetRMS()
	myMutex.Lock()
	Amount, err := strconv.ParseFloat(RMS.AvailableCash, 64)
	amount = Amount
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Second)
	myMutex.Unlock()
}

func TrackOrders(client *smartapigo.Client, symbol string) {
	for {
		//orders, _ := client.GetOrderBook()
		time.Sleep(1 * time.Second)
		positions, _ := client.GetPositions()
		isAnyPostionOpen := false
		totalPL := 0.0
		fmt.Printf("\nPositions %v\n", positions)

		for _, postion := range positions {
			qty, _ := strconv.Atoi(postion.NetQty)
			if postion.SymbolName == symbol && qty != 0 {
				pl, _ := strconv.ParseFloat(postion.NetValue, 64)
				fmt.Printf("current P/L in %v symbol is %v", symbol, pl)
			}
			if qty != 0 {
				isAnyPostionOpen = true
			}
			val, _ := strconv.ParseFloat(postion.NetValue, 64)
			totalPL += val
		}
		if isAnyPostionOpen == false {
			if totalPL <= -1000.0 || totalPL >= 2000.0 {
				CloseSession(client)
			}
			fmt.Printf("total P/L  %v", totalPL)
			setFlag(false)
			return
		}

	}

}

func setFlag(val bool) {
	myMutex.Lock()
	flag = val
	myMutex.Unlock()

}

func CalculatePosition(buyPrice, sl float64) int {
	myMutex.Lock()
	Amount := amount
	if Amount/buyPrice <= 1 {
		return 0
	}

	//maxRiskPercent := 0.05 // 2% maximum risk allowed
	//maxRiskAmount := Amount * maxRiskPercent
	//riskPerShare := math.Max(1, math.Abs(buyPrice-sl))
	//positionSize := maxRiskAmount / riskPerShare
	myMutex.Unlock()
	//return int(math.Min(Amount/buyPrice, positionSize))
	return int(Amount/buyPrice) * 4
}
