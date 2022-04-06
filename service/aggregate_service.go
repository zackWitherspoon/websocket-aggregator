package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/aggregate"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

const (
	aggregateTime = 5
)

type AggregateServicer interface {
	InitiateAggregateCalculation(tickerName string, conn *websocket.Conn, done chan bool)
	AddTradesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, done chan bool)
	SendTradesToChan(trades chan []byte, conn *websocket.Conn)
	ProcessTrades(tickerName string, tradesQueue chan []trade.TradeRequest, trades chan []byte, done chan bool)
}
type AggregateService struct{}

func (aggService *AggregateService) InitiateAggregateCalculation(tickerName string, conn *websocket.Conn, done chan bool) {
	trades := make(chan []byte, 10000)
	tradesQueue := make(chan []trade.TradeRequest, 10000)

	go aggService.ProcessTrades(tickerName, tradesQueue, trades, done)
	go aggService.AddTradesToBufferedChan(trades, tradesQueue, done)
	aggService.SendTradesToChan(trades, conn)
}

func (aggService *AggregateService) ProcessTrades(tickerName string, tradesQueue chan []trade.TradeRequest, trades chan []byte, done chan bool) {
	var tradesList []trade.TradeRequest
	lock := sync.Mutex{}
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
			tradesList = []trade.TradeRequest{}
			lock.Unlock()
		case <-done:
			return
		}
	}
}

func (aggService *AggregateService) SendTradesToChan(trades chan []byte, conn *websocket.Conn) {
	logrus.Debug("Scanner from connection -> []byte running")
	for {
		var msg []byte
		_, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		trades <- msg
	}
}

func (aggService *AggregateService) AddTradesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, done chan bool) {
	logrus.Debug("Scanner for []byte -> []trade.TradeRequest running")
	for {
		select {
		case <-done:
			return
		case msgBytes := <-trades:
			var m []trade.TradeRequest
			if err := json.Unmarshal(msgBytes, &m); err != nil {
				logrus.Debugf("An error was encountered with one of the incoming trades. That trade looks like:\n\t %s\n", string(msgBytes))
			} else {
				tradesQueue <- m
			}
		}
	}
}
