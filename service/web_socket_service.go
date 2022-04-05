package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/web_socket"
	"time"
)

const (
	webSocketUrl = "wss://delayed.polygon.io/stocks"
	APIKey       = "JncA2_tG82oqnPqGQHyMj2BW_h2jCoQl"
)

type WebSocketClient interface {
	InitializeConnection(ticker time.Ticker) *websocket.Conn
}

type TradeWebSocket struct{}

var authenticationMessage = []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKey))
var subscribeMessage = "{\"action\":\"subscribe\",\"params\":\"T.TSLA\"}"

func (tradeWS *TradeWebSocket) InitializeWSConnection(tickerName string) *websocket.Conn {
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
	return conn
}
