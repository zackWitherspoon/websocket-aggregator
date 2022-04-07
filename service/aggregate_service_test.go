package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"polygon-websocket-aggregator/model/trade"
	"time"
)

const tickerName = "TSLA"

//type mockAggregateService struct{}
//
//func (mockAggregateService) InitiateAggregateSequence(tickerName string, conn *websocket.Conn) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) AddTradeObjectsToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) AddIncomingBytesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, conn *websocket.Conn) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) ProcessTradesChan(tickerName string, tradesQueue chan []trade.TradeRequest, trades chan []byte) {
//	//TODO: RETURN SOMETHING HERE
//}
var twoTrades = []byte(`[{"ev":"T","sym":"TSLA","i":"15589","x":19,"p":1114.38,"s":13,"c":[14,37,41],"t":1649169893847,"q":1769466,"z":3},{"ev":"T","sym":"TSLA","i":"15590","x":19,"p":1114.37,"s":6,"c":[14,37,41],"t":1649169893847,"q":1769467,"z":3}]`)
var oneTrades = []byte(`[{"ev":"T","sym":"TSLA","i":"82186","x":4,"p":1114.31,"s":1,"c":[37],"t":1649169893984,"q":1769499,"z":3}]`)

func pushTradesToChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, done chan bool) {
	mockPushTradesToChan(trades, tradesQueue, done)
}

var mockPushTradesToChan func(trades chan []byte, tradesQueue chan []trade.TradeRequest, done chan bool)

var _ = Describe("Aggregate Test", func() {
	var (
		trades      chan []byte
		tradesQueue chan []trade.TradeRequest
		done        chan bool
		aggService  AggregateService
	)

	BeforeEach(func() {
		aggService = AggregateService{}
		trades = make(chan []byte, 10000)
		tradesQueue = make(chan []trade.TradeRequest, 10000)
		done = make(chan bool)

	})

	Describe("Given a connection has been established", func() {
		Context("When there are messages in the trades chan"+
			"And these messages have only 1 trade object in them", func() {
			It("Should be process all messages and change them into []Trades", func() {
				go aggService.AddTradeObjectsToBufferedChan(trades, tradesQueue, done)
				for i := 0; i < 10; i++ {
					trades <- oneTrades
					time.Sleep(500 * time.Millisecond)
					println("Length of tradesQueue =", len(tradesQueue))
				}
				println("Sending done")
				done <- true
				expected := 10
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When there are messages in the trades chan"+
			"And these messages have 2 trade object in them", func() {
			It("Should be process all messages and change them into []Trades", func() {
				go aggService.AddTradeObjectsToBufferedChan(trades, tradesQueue, done)
				for i := 0; i < 10; i++ {
					trades <- twoTrades
					time.Sleep(500 * time.Millisecond)
					println("Length of tradesQueue =", len(tradesQueue))
				}
				println("Sending done")
				done <- true
				expected := 10
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
