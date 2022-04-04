package aggregate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"polygon-websocket-aggregator/model/trade"
)

var _ = Describe("Aggregate Test", func() {
	arrayOfTrades := []trade.TradeRequest{
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 123.45, Size: 21, Conditions: []int{2, 12}, Timestamp: 1536036818784},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 111.23, Size: 11, Conditions: []int{1, 4}, Timestamp: 1536036818782},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 145.10, Size: 4, Conditions: []int{4, 1}, Timestamp: 1536036818786},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 90.26, Size: 25, Conditions: []int{12, 1}, Timestamp: 1536036818794},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 110.12, Size: 9, Conditions: []int{1}, Timestamp: 1536036818799},
	}
	startTimestamp := int64(1536036818784)
	Describe("Given a slice of Trades", func() {
		Context("When that slice of Trades is non-empty", func() {
			It("Should correctly match the expected aggregate", func() {

				expected := Aggregate{Symbol: "APPL", OpenPrice: 123.45, ClosingPrice: 110.12, HighPrice: 145.10, LowPrice: 90.26, Volume: 70, Timestamp: startTimestamp}
				actual := CalculateAggregate(arrayOfTrades, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that slice of Trades is empty", func() {
			It("Should correctly return an empty Aggregate", func() {
				expected := Aggregate{Symbol: "APPL", OpenPrice: 0, ClosingPrice: 0, HighPrice: 0, LowPrice: 0, Volume: 0, Timestamp: startTimestamp}
				actual := CalculateAggregate([]trade.TradeRequest{}, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that slice of Trades is empty", func() {
			It("Should correctly return an empty Aggregate", func() {
				expected := Aggregate{Symbol: "APPL", OpenPrice: 0, ClosingPrice: 0, HighPrice: 0, LowPrice: 0, Volume: 0, Timestamp: startTimestamp}
				actual := CalculateAggregate([]trade.TradeRequest{}, "APPL", startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
	})
	Describe("Given an Aggregate", func() {
		Context("When that Aggregate needs to be printed", func() {
			It("Should correctly match the expected aggregate", func() {
				agg := Aggregate{Symbol: "APPL", OpenPrice: 123.45, ClosingPrice: 145.22, HighPrice: 160.87, LowPrice: 111.23, Volume: 210, Timestamp: 1649022180}
				expected := "17:43:00 - open: $123.45, close: $145.22, high: $160.87, low: $111.23, volume: 210\n"
				actual := agg.PrintAggregate()
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
