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

type AggregateServicer interface {
	InitiateAggregateCalculation(tickerName string, tickerDuration time.Duration, conn *websocket.Conn, done chan bool)
	AddTradesToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, done chan bool)
	SendTradesToChan(trades chan []byte, conn *websocket.Conn)
	ProcessTrades(tickerName string, tradesQueue chan []trade.TradeRequest, done chan bool)
}
type AggregateService struct{}

func (aggService *AggregateService) InitiateAggregateCalculation(tickerName string, tickerDuration time.Duration, timeToKeepAggregates time.Duration, conn *websocket.Conn, done chan bool) {
	trades := make(chan []byte, 10000)
	tradesQueue := make(chan []trade.TradeRequest, 10000)

	go aggService.ProcessTrades(tickerName, tickerDuration, timeToKeepAggregates, tradesQueue, done)
	go aggService.AddTradesToBufferedChan(trades, tradesQueue, done)
	aggService.SendTradesToChan(trades, conn)
}

func (aggService *AggregateService) ProcessTrades(tickerName string, tickerDuration time.Duration, timeToKeepAggregates time.Duration, tradesQueue chan []trade.TradeRequest, done chan bool) {
	println("tickerDuration =", tickerDuration)
	println("timeToKeepAggregates =", timeToKeepAggregates)
	var tradesList []trade.TradeRequest
	var aggMap = make(map[aggregate.Duration]*aggregate.Aggregate)
	var startTime int64
	var currentSegmentTime int64
	timeHasElapsed := false
	lock := sync.Mutex{}
	aggMapLock := &sync.RWMutex{}
	ticker := time.NewTicker(tickerDuration)
	keepAggregates := time.NewTicker(timeToKeepAggregates)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			lock.Lock()
			startTime, currentSegmentTime, tradesList = aggService.updateAggMap(tickerName, tickerDuration, startTime, currentSegmentTime, tradesList, t, aggMap, timeHasElapsed)
			lock.Unlock()
		case tradeSlice := <-tradesQueue:
			for i := range tradeSlice {
				tradeSlice[i].PrintTrade()
				//TODO: Could improve this
				if currentSegmentTime == 0 {
					logrus.Debug("Setting current time & Start Time")
					currentSegmentTime = tradeSlice[i].Timestamp
					startTime = currentSegmentTime
					lock.Lock()
					tradesList = append(tradesList, tradeSlice[i])
					lock.Unlock()
					println("\t\t Just Set Start time to ", startTime)
					println("\t\t Just Set currentSegmentTime to ", currentSegmentTime)
				} else if tradeSlice[i].Timestamp > currentSegmentTime {
					logrus.Debug("Current Timestamp > currentSegmentTime")
					lock.Lock()
					tradesList = append(tradesList, tradeSlice[i])
					lock.Unlock()
				} else {
					println("Found a piece that was outside of the currentSegmentTime...", len(aggMap))
					if len(aggMap) != 0 {
						aggMapLock.RLock()
						aggMap = UpdatePastAgg(tickerName, aggMap, tradeSlice, i)
						aggMapLock.RUnlock()
					}
				}
			}
		case <-keepAggregates.C:
			aggMapLock.RLock()
			aggMap = PruneOldAggregates(aggMap, &startTime, timeToKeepAggregates)
			aggMapLock.RUnlock()
			keepAggregates.Stop()
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
				println("len(tadesQueue)=", len(tradesQueue))
				println("len(tadesQueue)=", len(trades))
				tradesQueue <- m
			}
		}
	}
}
