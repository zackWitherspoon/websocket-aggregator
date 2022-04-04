package aggregate

import (
	"fmt"
	"math"
	"polygon-websocket-aggregator/model/trade"
	"time"
)

type Aggregate struct {
	Symbol       string  `json:"sym"`
	OpenPrice    float64 `json:""`
	ClosingPrice float64 `json:""`
	HighPrice    float64 `json:""`
	LowPrice     float64 `json:""`
	Volume       int     `json:"v"`
	Timestamp    int64   `json:""`
}

func (agg *Aggregate) PrintAggregate() string {
	timestamp := time.Unix(agg.Timestamp, 0)
	return fmt.Sprintf("%d:%d:%.2d - open: $%.2f, close: $%.2f, high: $%.2f, low: $%.2f, volume: %d\n",
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(), agg.OpenPrice, agg.ClosingPrice, agg.HighPrice, agg.LowPrice, agg.Volume)
}

func CalculateAggregate(trades []trade.TradeRequest, symbol string, timeStamp int64) Aggregate {
	agg := Aggregate{Symbol: symbol, Volume: 0, Timestamp: timeStamp}
	var highestPrice float64
	var lowestPrice float64
	if len(trades) == 0 {
		return agg
	} else {
		agg.OpenPrice = trades[0].Price
		agg.ClosingPrice = trades[len(trades)-1].Price
		highestPrice = math.SmallestNonzeroFloat64
		lowestPrice = math.MaxFloat64
	}
	totalVolume := 0
	for _, t := range trades {
		totalVolume += t.Size
		if t.Price > highestPrice {
			highestPrice = t.Price
		}
		if t.Price < lowestPrice {
			lowestPrice = t.Price
		}
	}
	agg.LowPrice = lowestPrice
	agg.HighPrice = highestPrice
	agg.Volume = totalVolume
	return agg

}
