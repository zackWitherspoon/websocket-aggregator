package service

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"time"
)

func UpdatePastAgg(tickerName string, aggMap map[aggregate.Duration]*aggregate.Aggregate, tradeSlice []trade.TradeRequest, i int) map[aggregate.Duration]*aggregate.Aggregate {
	logrus.Debugf("Current Timestamp Outside the bounds of what we can keep track of. Searching through %d aggMap size", len(aggMap))
	for key, element := range aggMap {
		if key.Between(tradeSlice[i].Timestamp) {
			logrus.Debug("Current Timestamp is updating an older timeStamp")
			element.UpdateAggregate(tradeSlice[i], tickerName, tradeSlice[i].Timestamp)
			element.PrintAggregate()
		}
	}
	return aggMap
}

func (aggService *AggregateService) updateAggMap(tickerName string, tickerDuration time.Duration, startTime int64, currentSegmentTime int64, tradesList []trade.TradeRequest, t time.Time, aggMap map[aggregate.Duration]*aggregate.Aggregate, timeHasElapsed bool) (int64, int64, []trade.TradeRequest) {
	agg := aggregate.CalculateAggregate(tradesList, tickerName, t.Unix()-(time.Second*tickerDuration).Milliseconds())
	if len(aggMap) == 0 {
		key := aggregate.Duration{
			StartTime: agg.ClosingPriceTimestamp - tickerDuration.Milliseconds(),
			EndTime:   agg.ClosingPriceTimestamp,
		}
		aggMap[key] = &agg
		startTime = agg.ClosingPriceTimestamp - tickerDuration.Milliseconds()
		currentSegmentTime = agg.ClosingPriceTimestamp
	} else {
		key := aggregate.Duration{
			StartTime: currentSegmentTime,
			EndTime:   currentSegmentTime + tickerDuration.Milliseconds(),
		}
		aggMap[key] = &agg
		currentSegmentTime = currentSegmentTime + tickerDuration.Milliseconds()
		if timeHasElapsed {
			startTime += tickerDuration.Milliseconds()
		}
	}
	agg.PrintAggregate()
	tradesList = []trade.TradeRequest{}
	return startTime, currentSegmentTime, tradesList
}

func PruneOldAggregates(aggMap map[aggregate.Duration]*aggregate.Aggregate, startTime *int64, timeToKeepAggregates time.Duration) map[aggregate.Duration]*aggregate.Aggregate {
	for key, _ := range aggMap {
		if key.StartTime > *startTime+timeToKeepAggregates.Milliseconds() {
			a := timeToKeepAggregates.Milliseconds()
			startTime = &a
			delete(aggMap, key)
		}
	}
	return aggMap
}
