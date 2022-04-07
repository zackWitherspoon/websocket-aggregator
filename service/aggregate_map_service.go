package service

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

func UpdatePastAgg(aggMap map[aggregate.Duration]*aggregate.Aggregate, trade trade.TradeRequest) map[aggregate.Duration]*aggregate.Aggregate {
	logrus.Debugf("Current Timestamp Outside the bounds of what we can keep track of. Searching through %d aggMap size", len(aggMap))
	for aggDuration, aggregate := range aggMap {
		if aggDuration.Between(trade.Timestamp) {
			logrus.Debugf("\tUpdating Agg with timeStamp: %d and Volume: %d", aggregate.Timestamp, aggregate.Volume)
			logrus.Debugf("\tUpdating Agg with timeStamp: ", trade.Timestamp)
			aggregate.Update(trade)
			logrus.Debugf("\tPrinting Agg with timeStamp: %d and Volume: %d", trade.Timestamp, aggregate.Volume)
			aggregate.Print()
		}
	}
	return aggMap
}

func UpdateAggMap(tickerName string, tickerDuration time.Duration, rollingStartWindowTimestamp int64, rollingCurrentWindowTimestamp int64, tradesList []trade.TradeRequest, aggMap map[aggregate.Duration]*aggregate.Aggregate, rollingTimeWindowEnabled bool, aggMapLock *sync.RWMutex) (int64, int64, []trade.TradeRequest) {

	agg := aggregate.Calculate(tradesList, tickerName, rollingCurrentWindowTimestamp)
	key := aggregate.Duration{}
	if len(aggMap) == 0 {
		key.StartTime = agg.ClosingPriceTimestamp - tickerDuration.Milliseconds()
		key.EndTime = agg.ClosingPriceTimestamp
		rollingStartWindowTimestamp = agg.ClosingPriceTimestamp - tickerDuration.Milliseconds()
		rollingCurrentWindowTimestamp = agg.ClosingPriceTimestamp
	} else {
		key.StartTime = rollingCurrentWindowTimestamp
		key.EndTime = rollingCurrentWindowTimestamp + tickerDuration.Milliseconds()
		rollingCurrentWindowTimestamp = rollingCurrentWindowTimestamp + tickerDuration.Milliseconds()
		if rollingTimeWindowEnabled {
			aggMapLock.RLock()
			logrus.Debugf("Time to keep Aggregates has elipsed. There are currently %d items in the Aggregate Map\n", len(aggMap))
			aggMap = PruneExpiredAggregates(aggMap, rollingStartWindowTimestamp)
			logrus.Debugf("After pruning Aggregates, we are left with  %d items in the Aggregate Map\n", len(aggMap))
			aggMapLock.RUnlock()
			rollingStartWindowTimestamp += tickerDuration.Milliseconds()
		}
	}
	aggMap[key] = &agg
	agg.Print()
	tradesList = []trade.TradeRequest{}
	return rollingStartWindowTimestamp, rollingCurrentWindowTimestamp, tradesList
}

// PruneExpiredAggregates TODO: Might want to update this to a Cache instead. Purging the expired aggregates will be easier in te future
func PruneExpiredAggregates(aggMap map[aggregate.Duration]*aggregate.Aggregate, startTime int64) map[aggregate.Duration]*aggregate.Aggregate {
	for aggDuration, aggregate := range aggMap {
		if aggDuration.StartTime-startTime < 0 {
			logrus.Debugf("\t\tPruning an Aggregate from the Map: ")
			aggregate.DebugAggregate()
			delete(aggMap, aggDuration)
		}
	}
	return aggMap
}
