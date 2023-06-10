package main

import (
	"fmt"
	"os"

	SmartApi "github.com/TredingInGo/smartapi"
)

const (
	clientCode   = "P51284799"
	password     = "4926"
	apiKey       = "MN9K2rhC"
	marketKey    = "XDnby4up"
	totp         = "335589"
	refToken     = "eyJhbGciOiJIUzUxMiJ9.eyJ0b2tlbiI6IlJFRlJFU0gtVE9LRU4iLCJpYXQiOjE2ODQ3NzYwMjh9.l0RJix0Delk2kQyTBaIdtjngHzCx56R3g25BXoZ0LklJ5hbBXM2vkz_zRvDlt7pZCaVce8FYTC2xEaHRvz7SPQ"
	accessTokenC = "eyJhbGciOiJIUzUxMiJ9.eyJ1c2VybmFtZSI6IlA1MTI4NDc5OSIsInJvbGVzIjowLCJ1c2VydHlwZSI6IlVTRVIiLCJpYXQiOjE2ODQ3NzYwMjgsImV4cCI6MTY4NDg2MjQyOH0.fTwOimeaeEMOlPZCbsB45j6g8PmkBILu65o6JwfLeCJ7PTFjE-sTc7V_qT_yk90jzwhUsyY9QQvbuKqiHAAI6Q"
	feedTokenC   = "eyJhbGciOiJIUzUxMiJ9.eyJ1c2VybmFtZSI6IlA1MTI4NDc5OSIsImlhdCI6MTY4NDc3NjAyOCwiZXhwIjoxNjg0ODYyNDI4fQ.Nn5gchLf6p4YLM7-Q0PDKxlRuh4tT7Kjjl5LZRa5kBaPiukEIVrGD7rc9VLgvycMuZGZiAvzx749AVbojOdbqw"
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
	ABClient := SmartApi.New(clientCode, password, marketKey)
	//
	//fmt.Println("Client :- ", ABClient)
	//
	//User Login and Generate User Session
	if accessToken == "" {
		session, err := ABClient.GenerateSession(totp)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		setEnv(session)
	}

	//
	//////Renew User Tokens using refresh token
	////session.UserSessionTokens, err = ABClient.RenewAccessToken(session.RefreshToken)
	////
	////if err != nil {
	////	fmt.Println(err.Error())
	////	return
	////}
	//
	//fmt.Println("User Session Tokens :- ", session.UserSessionTokens)

	//Get User Profile
	//session.UserProfile, err = ABClient.GetUserProfile()

	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	//fmt.Println("User Profile :- ", session.UserProfile)
	//fmt.Println("User Session Object :- ", session)

	////Place Order
	//order, err := ABClient.PlaceOrder(SmartApi.OrderParams{Variety: "NORMAL", TradingSymbol: "SBIN-EQ", SymbolToken: "3045", TransactionType: "BUY", Exchange: "NSE", OrderType: "LIMIT", ProductType: "INTRADAY", Duration: "DAY", Price: "19500", SquareOff: "0", StopLoss: "0", Quantity: "1"})
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	//	fmt.Println("Placed Order ID and Script :- ", order)

	/*
			  "exchange": "NSE",
		     "symboltoken": "3045",
		     "interval": "ONE_MINUTE",
		     "fromdate": "2021-02-10 09:15",
		     "todate": "2021-02-10 09:16"
	*/

	//hc := clients.NewHistoricClient(clientCode, password)

	//session, err := hc.GenerateSession(totp)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	//hc.SetAccessToken(accessToken)
	//resp, err := hc.GetCandleData(clients.CandleParams{
	//	Exchange:    "NSE",
	//	SymbolToken: "3045",
	//	Interval:    "ONE_MINUTE",
	//	FromDate:    "2021-02-10 09:15",
	//	ToDate:      "2021-02-10 09:16",
	//})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//fmt.Println(resp)

	if accessToken != "" {
		ABClient.SetAccessToken(accessToken)
	}

	data, err := ABClient.GetCandleData(SmartApi.CandleParams{
		Exchange:    "NSE",
		SymbolToken: "3045",
		Interval:    "ONE_MINUTE",
		FromDate:    "2021-02-10 09:15",
		ToDate:      "2021-02-10 09:16",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	//
	fmt.Println(data)

	data, err = ABClient.GetCandleData(SmartApi.CandleParams{
		Exchange:    "NSE",
		SymbolToken: "3045",
		Interval:    "FIVE_MINUTE",
		FromDate:    "2023-02-10 09:15",
		ToDate:      "2023-02-10 09:21",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	//
	fmt.Println(data)
}

func setEnv(session SmartApi.UserSession) {
	os.Setenv("ACCESS_TOKEN", session.AccessToken)
	os.Setenv("FEED_TOKEN", session.FeedToken)
	os.Setenv("REFRESH_TOKEN", session.RefreshToken)
}
