package strategy

import (
	"database/sql"
	"fmt"
	"github.com/TredingInGo/AutomationService/historyData"
	smartapigo "github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
	"time"
)

type trade struct {
	spot      float64
	sl        float64
	tp        float64
	qty       float64
	orderType string
	flag      bool
}

type CandleResponse struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int
}

var tempOhlc = CandleResponse{
	Timestamp: time.Now(),
	Open:      0.0,
	High:      0.0,
	Low:       0.0,
	Close:     0.0,
	Volume:    0,
}

var t trade

const (
	startTime = "09:15"
	endTime   = "15:30"
)

var (
	currTime = time.Now()
	baseTime = time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 9, 0, 0, 0, time.Local)
)

type strategy struct {
	history  historyData.History
	pastData []smartapigo.CandleResponse

	LiveData    chan *models.SnapQuote
	chForCandle chan *models.SnapQuote
	db          *sql.DB
}

func New(history historyData.History, db *sql.DB) strategy {
	return strategy{
		history:     history,
		LiveData:    make(chan *models.SnapQuote, 100),
		chForCandle: make(chan *models.SnapQuote, 100),
		db:          db,
	}
}

func (s *strategy) Algo(token string) {
	fmt.Printf("Stock-- %v", GetStockName(token))
	go s.makeCandle(s.chForCandle, 300)

	for data := range s.LiveData {
		if len(s.pastData) == 0 {
			s.fillPastData(data.TokenInfo.Token)
		}

		s.chForCandle <- data
		// some algo....
		var closePrice = GetClosePriceArray(s.pastData)
		//closePrice = append(closePrice, float64(data.LastTradedPrice)/100.0)
		s.CalculateEma(closePrice, 9, token)
		s.CalculateSma(closePrice, 9, token)
		s.CalculateRsi(closePrice, 14, token)
		rsi := s.GetRsi(token)
		sma := s.GetSmaArray(token)
		ema := s.GetEmaArray(token)
		fmt.Printf(" sma = %v", sma[len(sma)-1])
		fmt.Printf(" ema = %v", ema[len(ema)-1])
		fmt.Printf(" Rsi = %v \n", rsi[len(rsi)-1])
		//closePrice = closePrice[:len(closePrice)-1]
		//s.oneRsStrategy(data)
		//s.Order(data)
		//fmt.Printf(" LiveData: ", float64(data.LastTradedPrice)/100.0)
	}
}
