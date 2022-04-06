package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/web_socket"
	"time"
)

const (
	APIKey = ""
)

type WebSocketClient interface {
	InitializeConnection(url string, ticker time.Ticker) *websocket.Conn
}

type TradeWebSocket struct{}

var authenticationMessage = []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKey))
var subscribeMessage = "{\"action\":\"subscribe\",\"params\":\"T.TSLA\"}"

func (tradeWS *TradeWebSocket) InitializeWSConnection(url string, tickerName string) *websocket.Conn {
	var responseMessage web_socket.WebSocketResponse

	logrus.Info("Attempting to connect to websocket at url: " + url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logrus.Fatalf("Dial to the websocket at URL: %s has failed. The following error was returned: %s\n", url, err.Error())
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
