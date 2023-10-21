package strategy

type stochField struct {
	k float64
	d float64
}

var rsi = make(map[string][]float64)
var atr = make(map[string][]float64)
var sma = make(map[string][]float64)
var ema = make(map[string][]float64)

func Init(stockName string) {
	ema[stockName] = make([]float64, 0)
	rsi[stockName] = make([]float64, 0)
	atr[stockName] = make([]float64, 0)
	sma[stockName] = make([]float64, 0)
}

// calculate ema
// multiplier = 2 / (period+1)
// ema[i] = (closePrice - ema[i-1) * multiplier + ema[i-1];

func (s *strategy) CalculateSma(data []float64, period int, stockName string) {
	smaArray := sma[stockName]
	if len(smaArray) == 0 {
		sum := 0.0
		for i := 0; i < period-1; i++ {
			smaArray = append(smaArray, -1.0)
			sum += data[i]
		}
		sum += data[period-1]
		smaArray = append(smaArray, sum/float64(period))
		for i := period; i < len(data); i++ {
			sum += data[i] - data[i-period]
			Current := sum / float64(period)
			smaArray = append(smaArray, Current)
		}
		sma[stockName] = smaArray
		return
	}
	if len(smaArray) < len(data) {
		sum := 0.0
		for j := len(smaArray) - period; j <= len(smaArray); j++ {
			sum += data[j]
		}
		smaArray = append(smaArray, sum/float64(period))
		for i := len(smaArray); i < len(data); i++ {
			sum += data[i] - data[i-period]
			Current := sum / float64(period)
			smaArray = append(smaArray, Current)
		}
	}
	sma[stockName] = smaArray
}
func (s *strategy) GetSmaArray(token string) []float64 {
	return sma[token]
}

func (s *strategy) CalculateEma(data []float64, period int, stockName string) {
	emaArray := ema[stockName]
	multiplier := 2.0 / float64(period+1)
	if len(emaArray) == 0 {
		sum := 0.0
		for i := 0; i < period-1; i++ {
			emaArray = append(emaArray, -1.0)
			sum += data[i]
		}
		sum += data[period-1]
		emaArray = append(emaArray, sum/float64(period))
		for i := period; i < len(data); i++ {
			Current := ((data[i] - emaArray[i-1]) * multiplier) + emaArray[i-1]
			emaArray = append(emaArray, Current)
		}
		ema[stockName] = emaArray
		return
	}
	if len(emaArray) < len(data) {
		for i := len(emaArray); i < len(data); i++ {
			Current := ((data[i] - emaArray[i-1]) * multiplier) + emaArray[i-1]
			emaArray = append(emaArray, Current)
		}
	}
	ema[stockName] = emaArray
}
func (s *strategy) GetEma(stockName string, ltp float64, period int) float64 {
	emaArray := ema[stockName]
	lastIdx := len(emaArray) - 1
	multiplier := 2.0 / float64(period+1)
	Current := ((ltp - emaArray[lastIdx]) * multiplier) + emaArray[lastIdx]
	return Current
}

func (s *strategy) GetEmaArray(stockName string) []float64 {
	emaArray := ema[stockName]

	return emaArray
}

func (s *strategy) CalculateRsi(data []float64, period int, stockName string) {
	rsiArray := rsi[stockName]

	var changeArray []float64
	var gainArray []float64
	var lossArray []float64
	for i := 0; i < len(data); i++ {
		if i == 0 {
			changeArray = append(changeArray, data[i])
		} else {
			changeArray = append(changeArray, data[i]-data[i-1])
		}

	}
	for i := 0; i < len(changeArray); i++ {
		if changeArray[i] >= 0 {
			gainArray = append(gainArray, changeArray[i])
			lossArray = append(lossArray, 0)
		} else {
			gainArray = append(gainArray, 0)
			lossArray = append(lossArray, -1*changeArray[i])
		}
	}
	stock1 := "Gain" + stockName
	stock2 := "Loss" + stockName
	s.CalculateSma(gainArray, period, stock1)
	s.CalculateSma(lossArray, period, stock2)
	emaGainArray := s.GetSmaArray(stock1)
	emaLossArray := s.GetSmaArray(stock2)
	if len(rsiArray) == 0 {
		for i := 0; i < len(data); i++ {
			avgGain := emaGainArray[i]
			avgLoss := emaLossArray[i]
			rs := avgGain / avgLoss
			rsiVal := 100 - (100 / (1 + rs))
			rsiArray = append(rsiArray, rsiVal)
		}
		rsi[stockName] = rsiArray
		return
	}
	if len(rsiArray) < len(data) {
		for i := len(rsiArray); i < len(data); i++ {
			avgGain := emaGainArray[i]
			avgLoss := emaLossArray[i]
			rs := avgGain / avgLoss
			rsiVal := 100.0 - (100.0 / (1 + rs))
			rsiArray = append(rsiArray, rsiVal)
		}
		rsi[stockName] = rsiArray
	}
}

func (s *strategy) GetRsi(stockName string) []float64 {
	return rsi[stockName]
}
