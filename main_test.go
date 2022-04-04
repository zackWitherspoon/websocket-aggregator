package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"polygon-websocket-aggregator/model/trade"
	"strings"
	"testing"
	"time"
)

var arrayOfTrades []trade.TradeRequest
var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	println("Size of array =", len(arrayOfTrades))
	for _, trade := range arrayOfTrades {
		a, _ := json.Marshal(trade)
		if err := c.WriteMessage(websocket.TextMessage, a); err != nil {
			println("%v", err)
		}

		if err != nil {
			println("ERROR!")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func TestExample(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	println(u)
	arrayOfTrades = []trade.TradeRequest{
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 123.45, Size: 21, Conditions: []int{2, 12}, Timestamp: 1536036818784},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 111.23, Size: 11, Conditions: []int{1, 4}, Timestamp: 1536036818782},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 145.10, Size: 4, Conditions: []int{4, 1}, Timestamp: 1536036818786},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 90.26, Size: 25, Conditions: []int{12, 1}, Timestamp: 1536036818794},
		trade.TradeRequest{Event: "T", Symbol: "APPL", ExchangeId: 1, TradeId: "12345", Tape: 1, Price: 110.12, Size: 9, Conditions: []int{1}, Timestamp: 1536036818799},
	}

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	//defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	chanMessages := make(chan trade.TradeRequest, 10000)

	// Read messages off the buffered queue:
	go func() {
		for msgBytes := range chanMessages {
			msgBytes.PrintTrade()
		}
	}()

	// As little logic as possible in the reader loop:
	for {
		var msg trade.TradeRequest
		ws.ReadJSON(&msg)
		if err != nil {
			println("ANOTHER ERROR!")
			println(err)
		}
		chanMessages <- msg
	}

	println("DONE")
}
