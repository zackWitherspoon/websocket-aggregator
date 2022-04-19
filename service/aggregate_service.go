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
	InitiateAggregateSequence(tickerName string, tickerDuration time.Duration, conn *websocket.Conn, testingInterruptChan chan bool)
	AddTradeObjectsToBufferedChan(trades chan []byte, tradesQueue chan []trade.TradeRequest, testingInterruptChan chan bool)
	AddIncomingBytesToBufferedChan(trades chan []byte, conn *websocket.Conn)
	ProcessTradesChan(tickerName string, tradesQueue chan []trade.TradeRequest, testingInterruptChan chan bool)
}
type AggregateService struct{}

func (aggService *AggregateService) InitiateAggregateSequence(tickerName string, tickerDuration time.Duration, aggregateCacheTime time.Duration, wsConn WebSocketClient, testingInterruptChan chan bool) {
	incomingByteChan := make(chan []byte, 10000)
	tradesListChan := make(chan []trade.TradeRequest, 10000)

	go aggService.ProcessTradesChan(tickerName, tickerDuration, aggregateCacheTime, tradesListChan, testingInterruptChan)
	go aggService.AddTradeObjectsToBufferedChan(incomingByteChan, tradesListChan, testingInterruptChan)
	aggService.AddIncomingBytesToBufferedChan(incomingByteChan, wsConn)
}

func (aggService *AggregateService) ProcessTradesChan(tickerName string, tickerDuration time.Duration, timeToKeepAggregates time.Duration, tradesQueue chan []trade.TradeRequest, testingInterruptChan chan bool) {

	var tradesList []trade.TradeRequest
	var aggMap = make(map[aggregate.Duration]*aggregate.Aggregate)
	var rollingStartWindowTimestamp, rollingCurrentWindowTimestamp int64
	rollingTimeWindowEnabled := false
	tradeProcessingLock := sync.Mutex{}
	aggMapLock := &sync.RWMutex{}
	aggregateWindow := time.NewTicker(tickerDuration)
	cacheEnabledTicker := time.NewTicker(timeToKeepAggregates)
	defer aggregateWindow.Stop()
	for {
		select {
		case <-aggregateWindow.C:
			logrus.Debugf("rollingStartWindowTimestamp before: %d & rollingCurrentWindowTimestamp: %dn", rollingStartWindowTimestamp, rollingCurrentWindowTimestamp)
			rollingStartWindowTimestamp, rollingCurrentWindowTimestamp = UpdateAggMap(tickerName, tickerDuration, rollingStartWindowTimestamp, rollingCurrentWindowTimestamp, tradesList, aggMap, rollingTimeWindowEnabled, aggMapLock)
			tradesList = []trade.TradeRequest{}
			logrus.Debugf("rollingStartWindowTimestamp after : %d & rollingCurrentWindowTimestamp :%d \n", rollingStartWindowTimestamp, rollingCurrentWindowTimestamp)
		case incomingTradeList := <-tradesQueue:
			for i := range incomingTradeList {
				if rollingCurrentWindowTimestamp == 0 {
					rollingCurrentWindowTimestamp = incomingTradeList[i].Timestamp
					rollingStartWindowTimestamp = rollingCurrentWindowTimestamp
					tradeProcessingLock.Lock()
					tradesList = append(tradesList, incomingTradeList[i])
					tradeProcessingLock.Unlock()
				} else if incomingTradeList[i].Timestamp > rollingCurrentWindowTimestamp {
					tradeProcessingLock.Lock()
					tradesList = append(tradesList, incomingTradeList[i])
					tradeProcessingLock.Unlock()
				} else {
					logrus.Debugf("Found a trade that was outside of the rollingCurrentWindowTimestamp: %d\n", len(aggMap))
					if len(aggMap) != 0 {
						aggMapLock.RLock()
						logrus.Debug("Calling UpdatePastAgg due to len aggMap not equaling 0")
						aggMap = UpdatePastAgg(aggMap, incomingTradeList[i])
						aggMapLock.RUnlock()
					}
				}
			}
		case <-cacheEnabledTicker.C:
			aggMapLock.RLock()
			logrus.Debugf("Time to keep Aggregates has elipsed. There are currently %d items in the Aggregate Map\n", len(aggMap))
			aggMap = PruneExpiredAggregates(aggMap, rollingStartWindowTimestamp)
			rollingTimeWindowEnabled = true
			logrus.Debugf("After pruning Aggregates, we are left with  %d items in the Aggregate Map\n", len(aggMap))
			aggMapLock.RUnlock()
			cacheEnabledTicker.Stop()
		case <-testingInterruptChan:
			return
		}
	}
}

func (aggService *AggregateService) AddIncomingBytesToBufferedChan(incomingByteChan chan []byte, wsConn WebSocketClient) {
	logrus.Debug("Scanner from connection -> []byte running")
	for {
		var msg []byte
		_, msg, err := wsConn.ReadMessage()
		//println(msg)
		if err != nil {
			panic(err)
		}
		incomingByteChan <- msg
	}
}

func (aggService *AggregateService) AddTradeObjectsToBufferedChan(incomingByteChan chan []byte, tradesListChan chan []trade.TradeRequest, testingInterruptChan chan bool) {
	logrus.Debug("Scanner for []byte -> []trade.TradeRequest running")
	for {
		select {
		case <-testingInterruptChan:
			return
		case msgBytes := <-incomingByteChan:
			var m []trade.TradeRequest
			if err := json.Unmarshal(msgBytes, &m); err != nil {
				logrus.Debugf("An error was encountered with one of the incoming incomingByteChan. That trade looks like:\n\t %s\n", string(msgBytes))
			} else {
				tradesListChan <- m
			}
		}
	}
}
