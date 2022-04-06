package integration_test

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"polygon-websocket-aggregator/model/trade"
	"polygon-websocket-aggregator/service"
	"strings"
	"time"
)

const tickerName = "TSLA"

type mockWebService struct{}

func (MWS mockWebService) InitializeWSConnection(tickerName string) *websocket.Conn {
	return createMockConn()
}

//type mockAggregateService struct{}
//
//func (mockAggregateService) InitiateAggregateCalculation(tickerName string, conn *websocket.Conn) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) AddTradesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) SendTradesToChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, conn *websocket.Conn) {
//	//TODO: RETURN SOMETHING HERE
//}
//func (mockAggregateService) ProcessTrades(tickerName string, tradesQueue chan []trade.TradeRequest, trades chan []byte) {
//	//TODO: RETURN SOMETHING HERE
//}
var twoTrades = []byte(`[{"ev":"T","sym":"TSLA","i":"15589","x":19,"p":1114.38,"s":13,"c":[14,37,41],"t":1649169893847,"q":1769466,"z":3},{"ev":"T","sym":"TSLA","i":"15590","x":19,"p":1114.37,"s":6,"c":[14,37,41],"t":1649169893847,"q":1769467,"z":3}]`)
var oneTrades = []byte(`[{"ev":"T","sym":"TSLA","i":"82186","x":4,"p":1114.31,"s":1,"c":[37],"t":1649169893984,"q":1769499,"z":3}]`)
var upgrader = websocket.Upgrader{}

var a []byte

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	for i := 1; i < 5; i++ {
		fmt.Printf("Adding the following message %d with text: %s\n", i, string(a))
		if err := c.WriteMessage(websocket.TextMessage, a); err != nil {
			println("%v", err)
		}
		if err != nil {
			println("ERROR!")
		}
		time.Sleep(1 * time.Second)
	}
}

func createMockConn() *websocket.Conn {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	println(u)
	// Connect to the server
	ws, _, _ := websocket.DefaultDialer.Dial(u, nil)
	return ws
}

var _ = Describe("Aggregate Test", func() {
	var (
		trades         chan []byte
		tradesQueue    chan []trade.TradeRequest
		done           chan bool
		mockAggService mockWebService
		aggService     service.AggregateService
	)

	BeforeEach(func() {
		mockAggService = mockWebService{}
		aggService = service.AggregateService{}
		trades = make(chan []byte, 10000)
		tradesQueue = make(chan []trade.TradeRequest, 10000)
		done = make(chan bool)

	})

	Describe("Given createMockConn connection to createMockConn websocket", func() {
		Context("When that websocket is sending createMockConn stream of tradeRequests", func() {
			It("Should correctly match the expected aggregate", func() {
				a = oneTrades
				mockConn := mockAggService.InitializeWSConnection(tickerName)
				go aggService.AddTradesToBufferedChan(trades, tradesQueue, done)
				go aggService.SendTradesToChan(trades, mockConn)
				time.Sleep(8 * time.Second)
				println("Sending done")
				done <- true
				expected := 4
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that websocket is sending createMockConn stream of tradeRequests", func() {
			It("Should correctly match the expected aggregate", func() {
				a = twoTrades
				mockConn := mockAggService.InitializeWSConnection(tickerName)
				go aggService.InitiateAggregateCalculation(tickerName, mockConn, done)
				time.Sleep(8 * time.Second)
				println("Sending done")
				done <- true
				expected := 0
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})

		//ADD Test to add []bytes to trades cahnnel
		//Context("When that websocket is sending createMockConn stream of tradeRequests", func() {
		//	It("Should correctly match the expected aggregate", func() {
		//		a = twoTrades
		//		mockConn := mockAggService.InitializeWSConnection(tickerName)
		//		go aggService.InitiateAggregateCalculation(tickerName, mockConn, done)
		//		time.Sleep(8 * time.Second)
		//		println("Sending done")
		//		done <- true
		//		expected := 0
		//		actual := len(tradesQueue)
		//		Expect(actual).To(Equal(expected))
		//	})
		//})
	})
})
