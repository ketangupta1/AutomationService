package main

import (
	"encoding/json"
	"fmt"
	"github.com/TredingInGo/AutomationService/Simulation"
	"github.com/TredingInGo/AutomationService/historyData"
	"github.com/TredingInGo/AutomationService/strategy"
	smartapi "github.com/TredingInGo/smartapi"
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

		m := map[string]string{}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Error reading request body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			http.Error(writer, "Error parsing JSON request body", http.StatusBadRequest)
			return
		}
		apiClient := smartapi.New(m["clientCode"], m["password"], m["marketKey"])
		session, err := apiClient.GenerateSession(m["totp"])
		if err != nil {
			errorMessage := fmt.Sprintf("Error generating session: %s", err.Error())
			http.Error(writer, errorMessage, http.StatusInternalServerError)
			return
		}
		setEnv(session)

		successMessage := fmt.Sprintf("User Session Tokens: %v", session.UserSessionTokens)
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{"message": "Connected successfully with angel one", "sessionTokens": successMessage})
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
		db := Simulation.Connect()
		strategy.TrendFollowingStretgy(apiClient, db)
	}).Methods(http.MethodPost)

	r.HandleFunc("/swing", func(writer http.ResponseWriter, request *http.Request) {
		if session.FeedToken == "" {
			fmt.Println("feed token not set")
			return
		}

		body, _ := ioutil.ReadAll(request.Body)
		//strategy.PopuletInstrumentsList()
		var param = make(map[string]string)
		json.Unmarshal(body, &param)
		db := Simulation.Connect()
		// Simulation.CollectData(db, apiClient) //this will populate list of stocks.
		strategy.SwingScreener(apiClient, db)

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
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8000"
	}

	http.ListenAndServe(":"+port, r)
}

func setEnv(session smartapi.UserSession) {
	os.Setenv("ACCESS_TOKEN", session.AccessToken)
	os.Setenv("FEED_TOKEN", session.FeedToken)
	os.Setenv("REFRESH_TOKEN", session.RefreshToken)

	feedToken = session.FeedToken
	accessToken = session.AccessToken
	refreshToken = session.RefreshToken
}
