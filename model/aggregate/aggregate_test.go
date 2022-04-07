package aggregate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"polygon-websocket-aggregator/model/trade"
	"sync"
)

var _ = Describe("Aggregate Test", func() {
	arrayOfTrades := []trade.TradeRequest{
		trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 123.45, Size: 21, Timestamp: 1536036818784},
		trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 111.23, Size: 11, Timestamp: 1536036818782},
		trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 145.10, Size: 4, Timestamp: 1536036818786},
		trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 90.26, Size: 25, Timestamp: 1536036818794},
		trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 110.12, Size: 9, Timestamp: 1536036818799},
	}
	startTimestamp := int64(1536036818784)
	var initialAggregate Aggregate

	BeforeEach(func() {
		initialAggregate = Aggregate{Symbol: "APPL", OpenPrice: 111.23, OpenPriceTimestamp: startTimestamp - 1, ClosingPrice: 110.12, ClosingPriceTimestamp: startTimestamp + 1, HighPrice: 145.10, LowPrice: 90.26, Volume: 70, Timestamp: startTimestamp, MutexLock: &sync.Mutex{}}
	})

	Describe("Given a slice of Trades", func() {
		Context("When that slice of Trades is non-empty", func() {
			It("Should correctly match the expected aggregate", func() {

				expected := Aggregate{Symbol: "APPL", OpenPrice: 111.23, OpenPriceTimestamp: 1536036818782, ClosingPrice: 110.12, ClosingPriceTimestamp: 1536036818799, HighPrice: 145.10, LowPrice: 90.26, Volume: 70, Timestamp: startTimestamp, MutexLock: &sync.Mutex{}}
				actual := Calculate(arrayOfTrades, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that slice of Trades is empty", func() {
			It("Should correctly return an empty Aggregate", func() {
				expected := Aggregate{Symbol: "APPL", OpenPrice: 0, ClosingPrice: 0, HighPrice: 0, LowPrice: 0, Volume: 0, Timestamp: startTimestamp, MutexLock: &sync.Mutex{}}
				actual := Calculate([]trade.TradeRequest{}, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that slice of Trades is empty", func() {
			It("Should correctly return an empty Aggregate", func() {
				expected := Aggregate{Symbol: "APPL", OpenPrice: 0, ClosingPrice: 0, HighPrice: 0, LowPrice: 0, Volume: 0, Timestamp: startTimestamp, MutexLock: &sync.Mutex{}}
				actual := Calculate([]trade.TradeRequest{}, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
	})
	Describe("Given a trade that has happened in the past", func() {
		Context("When the aggregate that represents the timespan that this trade is between exists"+
			"AND the trade has an Opening time that is before the current opening Time", func() {
			It("Should correctly update the aggregate", func() {

				incomingPastTrade := trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 125, Size: 21, Timestamp: startTimestamp - 4}
				expected := initialAggregate
				expected.OpenPrice = incomingPastTrade.Price
				expected.OpenPriceTimestamp = incomingPastTrade.Timestamp
				expected.Volume += incomingPastTrade.Size
				initialAggregate.Update(incomingPastTrade)
				Expect(initialAggregate).To(Equal(expected))
			})
		})
		Context("When the aggregate that represents the timespan that this trade is between exists"+
			"AND the trade has a closing time that is past the current Closing Time", func() {
			It("Should correctly update the aggregate", func() {

				incomingPastTrade := trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 125, Size: 21, Timestamp: startTimestamp + 3}
				expected := initialAggregate
				expected.ClosingPrice = incomingPastTrade.Price
				expected.ClosingPriceTimestamp = incomingPastTrade.Timestamp
				expected.Volume += incomingPastTrade.Size
				initialAggregate.Update(incomingPastTrade)
				Expect(initialAggregate).To(Equal(expected))
			})
		})
		Context("When the aggregate that represents the timespan that this trade is between exists"+
			"AND the trade has a high price than the high price", func() {
			It("Should correctly update the aggregate", func() {
				incomingPastTrade := trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 200.45, Size: 21, Timestamp: startTimestamp}
				expected := initialAggregate
				expected.HighPrice = incomingPastTrade.Price
				expected.Volume += incomingPastTrade.Size
				initialAggregate.Update(incomingPastTrade)
				Expect(initialAggregate).To(Equal(expected))
			})
		})
		Context("When the aggregate that represents the timespan that this trade is between exists"+
			"AND the trade has a lower price than the low price", func() {
			It("Should correctly update the aggregate", func() {
				incomingPastTrade := trade.TradeRequest{Symbol: "APPL", TradeId: "12345", Price: 2, Size: 21, Timestamp: startTimestamp}
				expected := initialAggregate
				expected.LowPrice = incomingPastTrade.Price
				expected.Volume += incomingPastTrade.Size
				initialAggregate.Update(incomingPastTrade)
				Expect(initialAggregate).To(Equal(expected))
			})
		})
	})
})
