package strategy

import (
	"fmt"
	"github.com/TredingInGo/AutomationService/clients"
	"github.com/TredingInGo/AutomationService/historyData"
	"github.com/TredingInGo/AutomationService/smartStream"
	smartapigo "github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
)

type strategy struct {
	history  historyData.History
	pastData []*clients.CandleResponse

	LiveData    chan *models.LTPInfo
	chForCandle chan *models.LTPInfo
}

func New(history historyData.History) strategy {
	return strategy{
		history:  history,
		LiveData: make(chan *models.LTPInfo, 100),
	}
}

func (s strategy) Algo() {
	go smartStream.MakeCandle(s.chForCandle, 60)

	for data := range s.LiveData {
		s.chForCandle <- data
		fmt.Println("LiveData: ", data)
	}
}

func (s strategy) fillPastData() {
	s.history.GetCandle(smartapigo.CandleParams{
		Exchange:    "",
		SymbolToken: "",
		Interval:    "",
		FromDate:    "",
		ToDate:      "",
	})
}
