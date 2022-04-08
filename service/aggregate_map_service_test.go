package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"sync"
)

var _ = Describe("Aggregate Test", func() {
	startTimestamp := int64(1536036818784)
	agg2Key := aggregate.Duration{StartTime: startTimestamp + 4*1000, EndTime: startTimestamp + 8*1000}
	var aggMap map[aggregate.Duration]*aggregate.Aggregate
	var agg1, agg2, agg3 aggregate.Aggregate

	BeforeEach(func() {
		agg1 = aggregate.Aggregate{Symbol: "APPL", OpenPrice: 111.23, OpenPriceTimestamp: startTimestamp, ClosingPrice: 110.12, ClosingPriceTimestamp: startTimestamp * 4 * 1000, HighPrice: 145.10, LowPrice: 90.26, Volume: 70, Timestamp: startTimestamp, MutexLock: &sync.Mutex{}}
		agg2 = aggregate.Aggregate{Symbol: "APPL", OpenPrice: 123.23, OpenPriceTimestamp: startTimestamp + 4*1000, ClosingPrice: 110.12, ClosingPriceTimestamp: startTimestamp + 8*10000, HighPrice: 156.10, LowPrice: 89.26, Volume: 11, Timestamp: startTimestamp + 4*1000, MutexLock: &sync.Mutex{}}
		agg3 = aggregate.Aggregate{Symbol: "APPL", OpenPrice: 134.23, OpenPriceTimestamp: startTimestamp + 8*1000, ClosingPrice: 110.12, ClosingPriceTimestamp: startTimestamp + 12*1000, HighPrice: 167.10, LowPrice: 91.26, Volume: 23, Timestamp: startTimestamp + 8*1000, MutexLock: &sync.Mutex{}}
		aggMap = map[aggregate.Duration]*aggregate.Aggregate{
			aggregate.Duration{StartTime: startTimestamp, EndTime: startTimestamp + 4*1000}:           &agg1,
			aggregate.Duration{StartTime: startTimestamp + 4*1000, EndTime: startTimestamp + 8*1000}:  &agg2,
			aggregate.Duration{StartTime: startTimestamp + 8*1000, EndTime: startTimestamp + 12*1000}: &agg3,
		}
	})

	Describe("Given a map of Aggregates", func() {
		Context("When an aggregate has been around for more than the expiration time", func() {
			It("Should prune the old aggregate", func() {
				startTime := startTimestamp + 4000
				expected := PruneExpiredAggregates(aggMap, startTime)
				Expect(len(expected)).To(Equal(2))
			})
		})
		Context("When a trade with an older timestamp comes in ", func() {
			It("Should update the aggregate for that time Duration", func() {
				tradeToUpdateMiddleAggregate := trade.TradeRequest{Symbol: "APPL", Price: 123.45, Size: 21, Timestamp: startTimestamp + 6*1000}
				for k, _ := range aggMap {
					println(k.StartTime)
					println(k.EndTime)
					println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
				}
				expectedMap := UpdatePastAgg(aggMap, tradeToUpdateMiddleAggregate)
				Expect(expectedMap[agg2Key].Volume).To(Equal(32))
			})
		})
	})
})
