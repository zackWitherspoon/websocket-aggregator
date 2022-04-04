package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"polygon-websocket-aggregator/model/web_socket"
	"time"
)

const (
	webSocketUrl = "wss://delayed.polygon.io/stocks"
	APIKey       = ""
)

type WebSocketClient interface {
	initializeConnection(ticker time.Ticker) (*websocket.Conn, error)
}

type TradeWebSocket struct{}

var authenticationMessage = []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKey))
var subscribeMessage = "{\"action\":\"subscribe\",\"params\":\"T.%s\"}"

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSocket := TradeWebSocket{}
	conn, err := webSocket.initializeWSConnection(tickerName)
	if err != nil {
		logrus.Fatalf("Websocket connection to the websocket at URL: %s has failed. The following error was returned: %s\n", webSocketUrl, err.Error())
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Printf("Tick at: %d", t)
			_, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("%v", err)
			}
			println(string(p))
		}
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
