package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"polygon-websocket-aggregator/model/web_socket"
	"sync"
	"time"
)

const (
	webSocketUrl  = "wss://delayed.polygon.io/stocks"
	APIKey        = "JncA2_tG82oqnPqGQHyMj2BW_h2jCoQl"
	aggregateTime = 5
)

type WebSocketClient interface {
	initializeConnection(ticker time.Ticker) (*websocket.Conn, error)
}

type TradeWebSocket struct{}

var authenticationMessage = []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKey))
var subscribeMessage = "{\"action\":\"subscribe\",\"params\":\"T.TSLA\"}"

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSocket := TradeWebSocket{}
	conn, err := webSocket.initializeWSConnection(tickerName)
	if err != nil {
		logrus.Fatalf("Websocket connection to the websocket at URL: %s has failed. The following error was returned: %s\n", webSocketUrl, err.Error())
	}
	defer conn.Close()
	//startTime := time.Now().Unix()
	ticker := time.NewTicker(aggregateTime * time.Second)
	defer ticker.Stop()
	trades := make(chan []byte, 10000)
	tradesQueue := make(chan []trade.TradeRequest, 10000)
	//aggMap := make(map[int64]*aggregate.Aggregate)
	var tradesList []trade.TradeRequest
	lock := sync.Mutex{}

	go func() {
		for {
			select {
			case t := <-tradesQueue:
				for i := range t {
					if !(t[i].Timestamp < time.Now().Unix()-time.Hour.Milliseconds()) {
						//println("Aquiring Lock for Update")
						lock.Lock()
						tradesList = append(tradesList, t[i])
						lock.Unlock()
						//println("Releasing Lock for Update")
					}
					//t[i].PrintTrade()
				}
			case t := <-ticker.C:
				//println("Aquiring Lock for Timer")
				lock.Lock()
				//agg := aggregate.Aggregate{}
				agg := aggregate.CalculateAggregate(tradesList, tickerName, t.Unix()-(time.Second*aggregateTime).Milliseconds())
				//println("~~~~~~~~~~~~~~~ABOUT TO PRINT AGGREGATE!")
				agg.PrintAggregate()
				println("Trades Length:", len(trades))
				println("TradesQueue Length:", len(trades))
				tradesList = []trade.TradeRequest{}
				lock.Unlock()
				//println("Releasing Lock for Timer")
			}

		}
	}()

	go func() {
		for msgBytes := range trades {
			var m []trade.TradeRequest
			if err := json.Unmarshal(msgBytes, &m); err != nil {
				panic(err)
			}
			tradesQueue <- m
			//logrus.Info("Message Bytes: ", msgBytes)
			// Send through Channel or add to slice.... you decide
		}
	}()
	for {

		var msg []byte
		_, msg, err := conn.ReadMessage()
		if err != nil {
			//TODO: Fix
			panic(err)
		}
		trades <- msg
	}
}

func (tradeWS *TradeWebSocket) initializeWSConnection(tickerName string) (*websocket.Conn, error) {
	var responseMessage web_socket.WebSocketResponse

	logrus.Info("Attempting to connect to websocket at url: " + webSocketUrl)
	conn, _, err := websocket.DefaultDialer.Dial(webSocketUrl, nil)
	if err != nil {
		logrus.Fatalf("Dial to the websocket at URL: %s has failed. The following error was returned: %s\n", webSocketUrl, err.Error())
	}

	err = conn.ReadJSON(&responseMessage)
	responseMessage.DebugResponse()
	//authenticate websocket
	authError := conn.WriteMessage(websocket.TextMessage, authenticationMessage)
	if authError != nil {
		logrus.Fatal(authError)
	}
	err = conn.ReadJSON(&responseMessage)
	responseMessage.DebugResponse()
	//subscribe to websocket
	var a = []byte(fmt.Sprintf(subscribeMessage, tickerName))
	subscribeError := conn.WriteMessage(websocket.TextMessage, a)
	err = conn.ReadJSON(&responseMessage)
	if err != nil {
		logrus.Fatal(subscribeError)
	}
	responseMessage.DebugResponse()
	return conn, nil
}
