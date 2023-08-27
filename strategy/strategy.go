package strategy

import (
	"fmt"
	"github.com/TredingInGo/AutomationService/historyData"
	smartapigo "github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
	"time"
)

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

	LiveData    chan *models.LTPInfo
	chForCandle chan *models.LTPInfo
}

func New(history historyData.History) strategy {
	return strategy{
		history:     history,
		LiveData:    make(chan *models.LTPInfo, 100),
		chForCandle: make(chan *models.LTPInfo, 100),
	}
}

func (s *strategy) Algo() {

	go s.makeCandle(s.chForCandle, 60)

	for data := range s.LiveData {
		if len(s.pastData) == 0 {
			// exchange, symbolToken, interval, startDate, endDate
			s.fillPastData(data.TokenInfo.Token)
		}

		s.chForCandle <- data
		fmt.Println("LiveData: ", data)
	}
}

func (s *strategy) fillPastData(symbol string) {
	// add last 10 days data
	var max = 10

	for count := 0; count < max; count++ {
		t := time.Now()
		t = t.Add(time.Hour * 24 * time.Duration(-count))

		if t.Weekday() == 0 || t.Weekday() == 6 {
			max++
			continue
		}

		year := t.Year()
		month := int(t.Month())
		day := t.Day()
		fromDate := fmt.Sprintf("%d-%02d-%02d %v", year, month, day, startTime)
		toDate := fmt.Sprintf("%d-%02d-%02d %v", year, month, day, endTime)

		candles, err := s.history.GetCandle(smartapigo.CandleParams{
			Exchange:    "NSE",
			SymbolToken: symbol,
			Interval:    "FIVE_MINUTE",
			FromDate:    fromDate,
			ToDate:      toDate,
		})
		if err != nil {
			fmt.Println("error while getting history data", err)
			continue
		}

		s.pastData = append(candles, s.pastData...)
	}

	fmt.Println(s.pastData, "\n\n")
}

func (s *strategy) makeCandle(ch <-chan *models.LTPInfo, duration int) {
	//candleDuration := 5

	//t, err := time.Parse("15:04:05", baseTime)
	//if err != nil {
	//	fmt.Println("Error parsing time:", err)
	//	return
	//}

	//formattedBaseTime := t.Format("15:04:05")
	lastSegStart := time.Time{}

	for data := range ch {
		epochSeconds := int64(data.ExchangeFeedTimeEpochMillis) / 1000
		dataTimeFormatted := time.Unix(epochSeconds, 0)
		if len(s.pastData) == 0 {
			s.pastData = append(s.pastData, smartapigo.CandleResponse{
				Timestamp: dataTimeFormatted,
				Open:      float64(data.LastTradedPrice) / 100,
				High:      float64(data.LastTradedPrice) / 100,
				Low:       float64(data.LastTradedPrice) / 100,
				Close:     float64(data.LastTradedPrice) / 100,
				Volume:    0,
			})
			tempTime := dataTimeFormatted.Sub(baseTime)
			fmt.Println("temp time", tempTime)
			tempTimeInSec := tempTime.Seconds()
			thresHoldTime := (int(tempTimeInSec)) / (duration)
			thresHoldTime++
			lastSegStart = baseTime.Add(time.Duration(thresHoldTime*duration) * time.Second)

		} else {
			if lastSegStart.After(dataTimeFormatted) {
				lastData := s.pastData[len(s.pastData)-1]
				ltp := float64(data.LastTradedPrice) / 100
				if lastData.Low > ltp {
					lastData.Low = ltp
				}

				if lastData.High < ltp {
					lastData.High = ltp
				}

				lastData.Close = ltp

				s.pastData[len(s.pastData)-1] = lastData
			} else {
				fmt.Println(s.pastData[len(s.pastData)-1])
				s.pastData = append(s.pastData, smartapigo.CandleResponse{
					Timestamp: dataTimeFormatted,
					Open:      float64(data.LastTradedPrice) / 100,
					High:      float64(data.LastTradedPrice) / 100,
					Low:       float64(data.LastTradedPrice) / 100,
					Close:     float64(data.LastTradedPrice) / 100,
					Volume:    0,
				})

				lastSegStart = lastSegStart.Add(time.Duration(duration) * time.Second)
			}
		}
	}

	for _, data := range s.pastData {
		fmt.Println(data)
	}
}
