package strategy

import (
	"fmt"
	"github.com/TredingInGo/AutomationService/clients"
	"github.com/TredingInGo/AutomationService/smartStream"
	"github.com/TredingInGo/smartapi/models"
)

type strategy struct {
	pastData []*clients.CandleResponse
	LiveData chan *models.LTPInfo
}

func New() strategy {
	return strategy{
		LiveData: make(chan *models.LTPInfo, 100),
	}
}

func (s strategy) Algo(takeHistory bool) {
	if takeHistory {
		fillPastData()
	}

	go smartStream.MakeCandle(s.LiveData, 60)

	for data := range s.LiveData {
		fmt.Println("LiveData: ", data)
	}
}

func fillPastData() {

}
