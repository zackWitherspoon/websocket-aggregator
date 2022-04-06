package aggregate

import (
	"fmt"
	"math"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

type Aggregate struct {
	Symbol                string  `json:"sym"`
	OpenPrice             float64 `json:""`
	OpenPriceTimestamp    int64
	ClosingPrice          float64 `json:""`
	ClosingPriceTimestamp int64
	HighPrice             float64 `json:""`
	LowPrice              float64 `json:""`
	Volume                int     `json:"v"`
	Timestamp             int64   `json:""`
	MutexLock             *sync.Mutex
}

func (agg *Aggregate) PrintAggregate() {
	timestamp := time.Unix(agg.Timestamp, 0)
	fmt.Printf("%d:%d:%.2d - open: $%.2f, close: $%.2f, high: $%.2f, low: $%.2f, volume: %d\n",
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(), agg.OpenPrice, agg.ClosingPrice, agg.HighPrice, agg.LowPrice, agg.Volume)
}

func (agg *Aggregate) UpdateAggregate(trade trade.TradeRequest, symbol string, timeStamp int64) {
	agg.MutexLock.Lock()
	defer agg.MutexLock.Unlock()
	if agg.Symbol == "" {
		agg = &Aggregate{Symbol: symbol, OpenPrice: trade.Price, OpenPriceTimestamp: trade.Timestamp, ClosingPrice: trade.Price,
			ClosingPriceTimestamp: trade.Timestamp, HighPrice: trade.Price, LowPrice: trade.Price, Volume: trade.Size,
			Timestamp: timeStamp}
		return
	}
	agg.Volume += trade.Size
	var highestPrice float64
	var lowestPrice float64
	if trade.Timestamp < agg.OpenPriceTimestamp {
		agg.OpenPrice = trade.Price
		agg.OpenPriceTimestamp = trade.Timestamp
	}
	if trade.Timestamp > agg.ClosingPriceTimestamp {
		agg.ClosingPrice = trade.Price
		agg.ClosingPriceTimestamp = trade.Timestamp
	}
	if trade.Price > highestPrice {
		highestPrice = trade.Price
	}
	if trade.Price < lowestPrice {
		lowestPrice = trade.Price
	}
}

func CalculateAggregate(trades []trade.TradeRequest, symbol string, timeStamp int64) Aggregate {
	agg := Aggregate{Symbol: symbol, Volume: 0, Timestamp: timeStamp}
	var highestPrice float64
	var lowestPrice float64
	if len(trades) == 0 {
		return agg
	} else {
		agg.OpenPriceTimestamp = trades[0].Timestamp
		agg.OpenPrice = trades[0].Price
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
		if t.Timestamp < agg.OpenPriceTimestamp {
			agg.OpenPrice = t.Price
			agg.OpenPriceTimestamp = t.Timestamp
		}
		if t.Timestamp > agg.ClosingPriceTimestamp {
			agg.ClosingPrice = t.Price
			agg.ClosingPriceTimestamp = t.Timestamp
		}
	}
	agg.LowPrice = lowestPrice
	agg.HighPrice = highestPrice
	agg.Volume = totalVolume
	return agg
}
