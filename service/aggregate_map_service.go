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
			logrus.Infof("\tUpdating Agg with timeStamp: %d and Volume: %d", element.Timestamp, element.Volume)
			logrus.Infof("\tUpdating Agg with timeStamp: ", tradeSlice.Timestamp)
			element.UpdateAggregate(tradeSlice)
			logrus.Infof("\tPrinting Agg with timeStamp: %d and Volume: %d", tradeSlice.Timestamp, element.Volume)
			element.PrintAggregate()
		}
	}
	return aggMap
}

func UpdateAggMap(tickerName string, tickerDuration time.Duration, timeToKeepAggregates time.Duration, startTime int64, currentSegmentTime int64,
	tradesList []trade.TradeRequest, t time.Time, aggMap map[aggregate.Duration]*aggregate.Aggregate,
	timeHasElapsed bool, aggMapLock *sync.RWMutex) (int64, int64, []trade.TradeRequest) {

	agg := aggregate.CalculateAggregate(tradesList, tickerName, currentSegmentTime)
	key := aggregate.Duration{}
	if len(aggMap) == 0 {
		key.StartTime = agg.ClosingPriceTimestamp - tickerDuration.Milliseconds()
		key.EndTime = agg.ClosingPriceTimestamp
		startTime = agg.ClosingPriceTimestamp - tickerDuration.Milliseconds()
		currentSegmentTime = agg.ClosingPriceTimestamp
	} else {
		key.StartTime = currentSegmentTime
		key.EndTime = currentSegmentTime + tickerDuration.Milliseconds()
		currentSegmentTime = currentSegmentTime + tickerDuration.Milliseconds()
		if timeHasElapsed {
			aggMapLock.RLock()
			logrus.Debugf("Time to keep Aggregates has elipsed. There are currently %d items in the Aggregate Map\n", len(aggMap))
			aggMap = PruneOldAggregates(aggMap, startTime, (timeToKeepAggregates).Milliseconds())
			logrus.Debugf("After pruning Aggregates, we are left with  %d items in the Aggregate Map\n", len(aggMap))
			aggMapLock.RUnlock()
			startTime += tickerDuration.Milliseconds()
		}
	}
	aggMap[key] = &agg
	agg.PrintAggregate()
	tradesList = []trade.TradeRequest{}
	return startTime, currentSegmentTime, tradesList
}

func PruneOldAggregates(aggMap map[aggregate.Duration]*aggregate.Aggregate, startTime int64, timeToKeepAggregates int64) map[aggregate.Duration]*aggregate.Aggregate {
	for key, element := range aggMap {
		if key.StartTime-startTime < 0 {
			logrus.Debugf("\t\tPruning an Aggregate from the Map: ")
			aggregate.DebugPrintAggregate(element)
			delete(aggMap, key)
		}
	}
	return aggMap
}
