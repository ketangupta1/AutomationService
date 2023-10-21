package strategy

import (
	"fmt"
	smartapigo "github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
	"sort"
	"time"
)

type params struct {
	key    string
	volume float64
}

type paramsSlice []params

func (s paramsSlice) Len() int {
	return len(s)
}

func (s paramsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s paramsSlice) Less(i, j int) bool {
	return s[i].key < s[j].key
}

var tokenMap = make(map[string][]smartapigo.CandleResponse)

func GetClosePriceArray(data []smartapigo.CandleResponse) []float64 {
	var closePrice []float64
	for i := 0; i < len(data); i++ {
		closePrice = append(closePrice, data[i].Close)
	}
	return closePrice
}

func GetVolumeArray(data []smartapigo.CandleResponse) []float64 {
	var volume []float64
	for i := 0; i < len(data); i++ {
		volume = append(volume, float64(data[i].Volume))
	}
	return volume
}

func (s strategy) FilterStocks(exchange string) []string {
	tokenList := GetAllToken(exchange)
	var filteredList []params
	for i := range tokenList {
		s.getPrevData(tokenList[i])
	}
	for key := range tokenMap {
		candles := tokenMap[key]
		lastCandle := candles[len(candles)-1]
		closingPrice := GetClosePriceArray(candles)
		volumes := GetVolumeArray(candles)
		last50 := closingPrice[len(closingPrice)-50:]
		last20Volume := volumes[len(volumes)-20:]
		s.CalculateSma(last20Volume, 9, key)
		s.CalculateEma(last50, 44, key)
		emaArray := s.GetEmaArray(key)
		smaArray := s.GetSmaArray(key)
		if lastCandle.Low > emaArray[len(emaArray)-1] && lastCandle.High >= 200.00 && lastCandle.High <= 1000.00 {
			filteredList = append(filteredList, params{
				key:    key,
				volume: smaArray[len(smaArray)-1],
			})
		}
	}
	sort.Sort(paramsSlice(filteredList))
	var keys []string
	for i := len(filteredList) - 3; i < len(filteredList); i++ {
		keys = append(keys, filteredList[i].key)
	}
	return keys
}

func (s strategy) getPrevData(token string) {
	pastData := tokenMap[token]
	var max = 3

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
			SymbolToken: token,
			Interval:    "FIVE_MINUTE",
			FromDate:    fromDate,
			ToDate:      toDate,
		})
		if err != nil {
			fmt.Println("error while getting history data", err)
			continue
		}

		pastData = append(candles, pastData...)
		tokenMap[token] = pastData
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
			Exchange:    "MCX",
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
}

func (s *strategy) makeCandle(ch <-chan *models.SnapQuote, duration int) {
	for data := range ch {
		epochSeconds := int64(data.ExchangeFeedTimeEpochMillis) / 1000
		dataTimeFormatted := time.Unix(epochSeconds, 0)
		lastCandleFormAt := s.pastData[len(s.pastData)-1].Timestamp
		nextCandleFormAt := lastCandleFormAt.Add(time.Second * time.Duration(duration))
		currentTime := time.Now()
		if currentTime.After(nextCandleFormAt) {
			s.pastData[len(s.pastData)-1].Timestamp = nextCandleFormAt
			tempOhlc = CandleResponse{
				Timestamp: dataTimeFormatted,
				Open:      float64(data.LastTradedPrice) / 100,
				High:      float64(data.LastTradedPrice) / 100,
				Low:       float64(data.LastTradedPrice) / 100,
				Close:     float64(data.LastTradedPrice) / 100,
				Volume:    0,
			}
			s.pastData = append(s.pastData, smartapigo.CandleResponse{
				Timestamp: tempOhlc.Timestamp,
				Open:      tempOhlc.Open,
				High:      tempOhlc.High,
				Low:       tempOhlc.Low,
				Close:     tempOhlc.Close,
				Volume:    0,
			})
		} else {
			if tempOhlc.Open == 0.0 {

				tempOhlc = CandleResponse{
					Timestamp: lastCandleFormAt,
					Open:      float64(data.LastTradedPrice) / 100,
					High:      float64(data.LastTradedPrice) / 100,
					Low:       float64(data.LastTradedPrice) / 100,
					Close:     float64(data.LastTradedPrice) / 100,
					Volume:    0,
				}
				s.pastData = append(s.pastData, smartapigo.CandleResponse{
					Timestamp: tempOhlc.Timestamp,
					Open:      tempOhlc.Open,
					High:      tempOhlc.High,
					Low:       tempOhlc.Low,
					Close:     tempOhlc.Close,
					Volume:    0,
				})
			} else {
				tempOhlc.Close = float64(data.LastTradedPrice) / 100
				s.pastData[len(s.pastData)-1].Close = tempOhlc.Close
			}

		}
	}

}
