package service

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

func UpdatePastAgg(aggMap map[aggregate.Duration]*aggregate.Aggregate, tradeSlice trade.TradeRequest) map[aggregate.Duration]*aggregate.Aggregate {
	logrus.Debugf("Current Timestamp Outside the bounds of what we can keep track of. Searching through %d aggMap size", len(aggMap))
	for key, element := range aggMap {
		if key.Between(tradeSlice.Timestamp) {
			logrus.Debug("Current Timestamp is updating an older timeStamp")
			element.UpdateAggregate(tradeSlice)
			println("\t\tPrinting Aggregate From the past")
			element.PrintAggregate()
		}
	}
	return aggMap
}

func UpdateAggMap(tickerName string, tickerDuration time.Duration, timeToKeepAggregates time.Duration, startTime int64, currentSegmentTime int64,
	tradesList []trade.TradeRequest, t time.Time, aggMap map[aggregate.Duration]*aggregate.Aggregate,
	timeHasElapsed bool, aggMapLock *sync.RWMutex) (int64, int64, []trade.TradeRequest) {

	agg := aggregate.CalculateAggregate(tradesList, tickerName, t.Unix()-(time.Second*tickerDuration).Nanoseconds()*1000)
	key := aggregate.Duration{}
	if len(aggMap) == 0 {
		key.StartTime = agg.ClosingPriceTimestamp - tickerDuration.Nanoseconds()
		key.EndTime = agg.ClosingPriceTimestamp
		startTime = agg.ClosingPriceTimestamp - tickerDuration.Nanoseconds()
		currentSegmentTime = agg.ClosingPriceTimestamp
	} else {
		key.StartTime = currentSegmentTime
		key.EndTime = currentSegmentTime + tickerDuration.Nanoseconds()
		currentSegmentTime = currentSegmentTime + tickerDuration.Nanoseconds()
		if timeHasElapsed {
			aggMapLock.RLock()
			aggMap = PruneOldAggregates(aggMap, &startTime, (timeToKeepAggregates).Nanoseconds())
			aggMapLock.RUnlock()
			startTime += tickerDuration.Nanoseconds()
		}
	}
	aggMap[key] = &agg
	println("\t\tPrinting Aggregate from TICKER")
	agg.PrintAggregate()
	tradesList = []trade.TradeRequest{}
	return startTime, currentSegmentTime, tradesList
}

func PruneOldAggregates(aggMap map[aggregate.Duration]*aggregate.Aggregate, startTime *int64, timeToKeepAggregates int64) map[aggregate.Duration]*aggregate.Aggregate {
	for key, _ := range aggMap {
		if key.StartTime > *startTime+timeToKeepAggregates {
			a := timeToKeepAggregates
			startTime = &a
			delete(aggMap, key)
		}
	}
	return aggMap
}
