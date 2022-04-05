package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

const (
	aggregateTime = 5
)

func InitiateAggregateCalculation(tickerName string, conn *websocket.Conn) {
	trades := make(chan []byte, 10000)
	tradesQueue := make(chan []trade.TradeRequest, 10000)
	go processTrades(tickerName, tradesQueue, trades)
	go addTradesToBufferedChan(trades, tradesQueue)
	sendTradesToChan(trades, tradesQueue, conn)
}

func processTrades(tickerName string, tradesQueue chan []trade.TradeRequest, trades chan []byte) {
	var tradesList []trade.TradeRequest
	lock := sync.Mutex{}
	count := 0
	println("new timestamp       start time")
	ticker := time.NewTicker(aggregateTime * time.Second)
	defer ticker.Stop()
	for {
		select {
		case tradeSlice := <-tradesQueue:
			for i := range tradeSlice {
				if !(tradeSlice[i].Timestamp < time.Now().Unix()-time.Hour.Milliseconds()) {
					lock.Lock()
					tradesList = append(tradesList, tradeSlice[i])
					lock.Unlock()
				}
			}
		case t := <-ticker.C:
			lock.Lock()
			agg := aggregate.CalculateAggregate(tradesList, tickerName, t.Unix()-(time.Second*aggregateTime).Milliseconds())
			agg.PrintAggregate()
			println("Trades Length:", len(trades))
			println("TradesQueue Length:", len(trades))
			println("Count currently at: ", count)
			count = 0
			tradesList = []trade.TradeRequest{}
			lock.Unlock()
		}
	}
}

func sendTradesToChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, conn *websocket.Conn) {
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

func addTradesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest) {
	for msgBytes := range trades {
		var m []trade.TradeRequest
		if err := json.Unmarshal(msgBytes, &m); err != nil {
			panic(err)
		}
		tradesQueue <- m
	}
}
