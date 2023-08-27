package main

import (
	"encoding/json"
	"fmt"
	"github.com/TredingInGo/AutomationService/historyData"
	"github.com/TredingInGo/AutomationService/smartStream"
	"github.com/TredingInGo/AutomationService/strategy"
	smartapi "github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
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
	apiClient                            *smartapi.Client
	session                              smartapi.UserSession
	err                                  error
)

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	feedToken = os.Getenv("FEED_TOKEN")
	refreshToken = os.Getenv("REFRESH_TOKEN")
}

func main() {
	defer func() {
		recover()
	}()

	r := mux.NewRouter()

	r.HandleFunc("/session", func(writer http.ResponseWriter, request *http.Request) {
		//// Create New Angel Broking Client
		apiClient = smartapi.New(clientCode, password, marketKey)

		m := map[string]string{}
		body, err := ioutil.ReadAll(request.Body)

		json.Unmarshal(body, &m)
		//User Login and Generate User Session
		session, err = apiClient.GenerateSession(m["totp"])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("User Session Tokens :- ", session.UserSessionTokens)
		setEnv(session)
	}).Methods(http.MethodPost)

	r.HandleFunc("/candle", func(writer http.ResponseWriter, request *http.Request) {
		history := historyData.New(apiClient)
		params := request.URL.Query()

		data, err := history.GetCandle(smartapi.CandleParams{
			Exchange:    params.Get("exchange"),
			SymbolToken: params.Get("symbolToken"),
			Interval:    params.Get("interval"),
			FromDate:    params.Get("fromDate"),
			ToDate:      params.Get("toDate"),
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(data)

		b, _ := json.Marshal(data)
		writer.Write(b)
		writer.WriteHeader(200)
	}).Methods(http.MethodGet)

	r.HandleFunc("/startStream", func(writer http.ResponseWriter, request *http.Request) {
		if session.FeedToken == "" {
			fmt.Println("feed token not set")
			return
		}

		body, _ := ioutil.ReadAll(request.Body)
		var param = make(map[string]string)
		json.Unmarshal(body, &param)

		history := historyData.New(apiClient)
		someAlgo := strategy.New(history)
		client := smartStream.New(clientCode, feedToken)
		go client.Connect(someAlgo.LiveData, models.LTP, models.NSECM, param["token"])

		go someAlgo.Algo()

	}).Methods(http.MethodPost)

	r.HandleFunc("/renew", func(writer http.ResponseWriter, request *http.Request) {
		//Renew User Tokens using refresh token
		session.UserSessionTokens, err = apiClient.RenewAccessToken(session.RefreshToken)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("User Session Tokens :- ", session.UserSessionTokens)
	}).Methods(http.MethodGet)

	r.HandleFunc("/profile", func(writer http.ResponseWriter, request *http.Request) {
		// Get User Profile
		session.UserProfile, err = apiClient.GetUserProfile()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("User Profile :- ", session.UserProfile)
		fmt.Println("User Session Object :- ", session)

	})

	r.HandleFunc("/order", func(writer http.ResponseWriter, request *http.Request) {
		//Place Order
		order, err := apiClient.PlaceOrder(smartapi.OrderParams{
			Variety:         "NORMAL",
			TradingSymbol:   "SBIN-EQ",
			SymbolToken:     "3045",
			TransactionType: "BUY",
			Exchange:        "NSE",
			OrderType:       "LIMIT",
			ProductType:     "INTRADAY",
			Duration:        "DAY",
			Price:           "19500",
			SquareOff:       "0",
			StopLoss:        "0",
			Quantity:        "1",
		})

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Placed Order ID and Script :- ", order)
	})

	http.ListenAndServe(":8000", r)
}

func setEnv(session smartapi.UserSession) {
	os.Setenv("ACCESS_TOKEN", session.AccessToken)
	os.Setenv("FEED_TOKEN", session.FeedToken)
	os.Setenv("REFRESH_TOKEN", session.RefreshToken)

	feedToken = session.FeedToken
	accessToken = session.AccessToken
	refreshToken = session.RefreshToken
}
