package main

import (
	"fmt"
	simulation "github.com/TredingInGo/AutomationService/Simulation"
	"github.com/TredingInGo/AutomationService/historyData"
	"github.com/TredingInGo/AutomationService/smartStream"
	"os"
	"sync"

	smartapi "github.com/TredingInGo/smartapi"
)

const (
	clientCode = "P51284799"
	password   = "4926"
	apiKey     = "MN9K2rhC"
	marketKey  = "XDnby4up"
	totp       = "874294"
)

var (
	accessToken, feedToken, refreshToken string
)

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	feedToken = os.Getenv("FEED_TOKEN")
	refreshToken = os.Getenv("REFRESH_TOKEN")
}

func main() {

	//// Create New Angel Broking Client
	apiClient := smartapi.New(clientCode, password, marketKey)

	//User Login and Generate User Session
	session, err := apiClient.GenerateSession(totp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	setEnv(session)

	wg := sync.WaitGroup{}

	client := smartStream.New(clientCode, feedToken)
	history := historyData.New(apiClient)
	//someAlgo := strategy.New(history)
	_ = client
	//go func() {
	//	wg.Add(1)
	//	defer wg.Done()
	//	client.Connect(someAlgo.LiveData, models.LTP, models.NSECM, "2885")
	//}()
	db := simulation.Connect()
	//simulation.PrepareData(db, apiClient)
	simulation.DoBackTest(db)
	//go someAlgo.Algo()

	////
	////////Renew User Tokens using refresh token
	//////session.UserSessionTokens, err = ABClient.RenewAccessToken(session.RefreshToken)
	//////
	//////if err != nil {
	//////	fmt.Println(err.Error())
	//////	return
	//////}
	////
	////fmt.Println("User Session Tokens :- ", session.UserSessionTokens)
	//
	////Get User Profile
	////session.UserProfile, err = ABClient.GetUserProfile()
	//
	////if err != nil {
	////	fmt.Println(err.Error())
	////	return
	////}
	//
	////fmt.Println("User Profile :- ", session.UserProfile)
	////fmt.Println("User Session Object :- ", session)
	//
	//////Place Order
	////order, err := ABClient.PlaceOrder(SmartApi.OrderParams{Variety: "NORMAL", TradingSymbol: "SBIN-EQ", SymbolToken: "3045", TransactionType: "BUY", Exchange: "NSE", OrderType: "LIMIT", ProductType: "INTRADAY", Duration: "DAY", Price: "19500", SquareOff: "0", StopLoss: "0", Quantity: "1"})
	////
	////if err != nil {
	////	fmt.Println(err.Error())
	////	return
	////}
	//
	////	fmt.Println("Placed Order ID and Script :- ", order)
	//
	///*
	//		  "exchange": "NSE",
	//	     "symboltoken": "3045",
	//	     "interval": "ONE_MINUTE",
	//	     "fromdate": "2021-02-10 09:15",
	//	     "todate": "2021-02-10 09:16"
	//*/
	//
	////hc := clients.NewHistoricClient(clientCode, password)
	//
	////session, err := hc.GenerateSession(totp)
	////if err != nil {
	////	fmt.Println(err.Error())
	////}
	//
	////hc.SetAccessToken(accessToken)
	////resp, err := hc.GetCandleData(clients.CandleParams{
	////	Exchange:    "NSE",
	////	SymbolToken: "3045",
	////	Interval:    "ONE_MINUTE",
	////	FromDate:    "2021-02-10 09:15",
	////	ToDate:      "2021-02-10 09:16",
	////})
	////if err != nil {
	////	fmt.Println(err.Error())
	////}
	////
	////fmt.Println(resp)

	data, err := history.GetCandle(smartapi.CandleParams{
		Exchange:    "NSE",
		SymbolToken: "3045",
		Interval:    "ONE_MINUTE",
		FromDate:    "2021-02-10 09:15",
		ToDate:      "2021-02-10 09:16",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(data)

	data, err = history.GetCandle(smartapi.CandleParams{
		Exchange:    "NSE",
		SymbolToken: "3045",
		Interval:    "FIVE_MINUTE",
		FromDate:    "2023-02-10 09:15",
		ToDate:      "2023-01-10 09:21",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	//
	fmt.Println(data)

	wg.Wait()
}

func setEnv(session smartapi.UserSession) {
	os.Setenv("ACCESS_TOKEN", session.AccessToken)
	os.Setenv("FEED_TOKEN", session.FeedToken)
	os.Setenv("REFRESH_TOKEN", session.RefreshToken)

	feedToken = session.FeedToken
	accessToken = session.AccessToken
	refreshToken = session.RefreshToken
}
