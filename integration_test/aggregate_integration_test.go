package integration_test

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"polygon-websocket-aggregator/application"
	"polygon-websocket-aggregator/model/trade"
	"polygon-websocket-aggregator/service"
	"strings"
	"time"
)

const tickerName = "TSLA"

type mockWebService struct{}

func (MWS mockWebService) InitializeWSConnection(string) *websocket.Conn {
	return createMockConn()
}

var outOfOrder = false
var upgrader = websocket.Upgrader{}

func writeTrades(w http.ResponseWriter, r *http.Request) {
	file := "6_trades.txt"
	if outOfOrder {
		file = "18_trades.txt"
	}
	c, err := upgrader.Upgrade(w, r, nil)
	tradesText, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer tradesText.Close()
	reader := bufio.NewScanner(tradesText)
	for reader.Scan() {
		fmt.Printf("Adding the following message with text: %s\n", string(reader.Text()))
		if err := c.WriteMessage(websocket.TextMessage, []byte(reader.Text())); err != nil {
			println("ERROR!!!! %v", err)
		}
		if err != nil {
			log.Fatalf("An error has occurred in writing the messages out %s \n", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
	if err := reader.Err(); err != nil {
		log.Fatalf("An error has occurred with the reader: %s\n", err.Error())
	}
}

func createMockConn() *websocket.Conn {
	s := httptest.NewServer(http.HandlerFunc(writeTrades))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	println(u)
	// Connect to the server
	ws, _, _ := websocket.DefaultDialer.Dial(u, nil)
	return ws
}

var (
	trades         chan []byte
	tradesQueue    chan []trade.TradeRequest
	done           chan bool
	mockAggService mockWebService
	aggService     service.AggregateService
)

var _ = Describe("Aggregate Test", func() {

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
				mockConn := mockAggService.InitializeWSConnection(tickerName)
				go aggService.AddTradeObjectsToBufferedChan(trades, tradesQueue, done)
				go aggService.AddIncomingBytesToBufferedChan(trades, mockConn)
				time.Sleep(10 * time.Second)
				fmt.Println("Sending done")
				done <- true
				expected := 6
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that websocket is sending createMockConn stream of tradeRequests", func() {
			It("Should correctly match the expected aggregate", func() {
				outOfOrder = true
				mockConn := mockAggService.InitializeWSConnection(tickerName)
				go aggService.InitiateAggregateSequence(tickerName, 5*application.SecondsInMilliseconds, 10*application.SecondsInMilliseconds, mockConn, done)
				time.Sleep(20 * time.Second)
				fmt.Println("Sending done")
				done <- true
				expected := 0
				actual := len(tradesQueue)
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
