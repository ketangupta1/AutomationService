package Simulation

import (
	"math"
	"time"
)

func CalculateSMA(data []OHLC, period int) []float64 {
	sma := make([]float64, len(data)+1)

	for i := 100; i < len(data); i++ {
		sum := 0.0
		for j := i; j > i-period; j-- {
			sum += data[j].Close
		}
		sma[i] = sum / float64(period)
	}

	return sma
}

func CalculateEMA(data []OHLC, period int) []float64 {
	ema := make([]float64, len(data)-period+1)

	sma := CalculateSMA(data, period)

	for i := 20; i < period; i++ {
		ema[i] = sma[i]
	}

	smoothingFactor := 2.0 / (float64(period) + 1)
	for i := period; i < len(data); i++ {
		ema[i-period+1] = (data[i].Close-sma[i-period])*smoothingFactor + ema[i-period]
	}

	return ema
}

func CalculateRSI(data []OHLC, index int, period int) float64 {
	rsiPeriod := period
	if index < rsiPeriod {
		return 0.0 // RSI not available for the first few data points
	}

	gain := 0.0
	loss := 0.0

	// Calculate average gain and loss
	for i := index; i > index-rsiPeriod; i-- {
		diff := data[i].Close - data[i-1].Close
		if diff >= 0 {
			gain += diff
		} else {
			loss -= diff
		}
	}
	avgGain := gain / float64(rsiPeriod)
	avgLoss := loss / float64(rsiPeriod)

	// Calculate RSI
	if avgLoss == 0 {
		return 100.0 // Avoid division by zero
	} else {
		rs := avgGain / avgLoss
		return 100.0 - (100.0 / (1.0 + rs))
	}
}

func CalculateBollingerBands(data []OHLC, period int, k float64) ([]float64, []float64) {
	sma := CalculateSMA(data, period)
	sd := calculateSD(data, period)
	ubb := make([]float64, len(data)-period+1)
	lbb := make([]float64, len(data)-period+1)

	for i := period; i < len(data); i++ {
		ubb[i-period] = sma[i-period] + (sd[i-period] * k)
		lbb[i-period] = sma[i-period] - (sd[i-period] * k)
	}

	return ubb, lbb
}

func calculateSD(data []OHLC, period int) []float64 {
	sd := make([]float64, len(data)-period)

	for i := period; i < len(data)-period; i++ {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += (data[j].Close - data[i-1].Close) * (data[j].Close - data[i-1].Close)
		}
		mean := sum / float64(period)
		sd[i-period] = math.Sqrt(mean)
	}
	return sd
}

type StochasticData struct {
	K float64
	D float64
}

func parseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04", timestamp)
}

func CalculateStochastic(ohlcData []OHLC, currentIndex int, period int) StochasticData {
	if currentIndex < period-1 {
		return StochasticData{}
	}

	var highestHigh, lowestLow float64

	for i := currentIndex - (period - 1); i <= currentIndex; i++ {
		if i == currentIndex-(period-1) {
			highestHigh = ohlcData[i].High
			lowestLow = ohlcData[i].Low
		} else {
			if ohlcData[i].High > highestHigh {
				highestHigh = ohlcData[i].High
			}
			if ohlcData[i].Low < lowestLow {
				lowestLow = ohlcData[i].Low
			}
		}
	}

	currentClose := ohlcData[currentIndex].Close
	k := (currentClose - lowestLow) / (highestHigh - lowestLow) * 100.0

	// Calculate D using a simple moving average
	d := calculateSMAFor(ohlcData, currentIndex, period)

	return StochasticData{K: k, D: d}
}

func calculateSMAFor(ohlcData []OHLC, currentIndex int, period int) float64 {
	sum := 0.0
	for i := currentIndex - (period - 1); i <= currentIndex; i++ {
		sum += ohlcData[i].Close
	}
	return sum / float64(period)
}
