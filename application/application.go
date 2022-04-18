package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
	"time"
)

const (
	//TODO: Before going to production, the following would want to be things we might want to grab from env variables/volume
	webSocketUrl              = "wss://delayed.polygon.io/stocks"
	aggregateTimeWindow       = 30 * time.Second
	aggregatePersistentWindow = 3600 * time.Second
)

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSocket := service.TradeWebSocket{}
	wsConn := webSocket.InitializeWSConnection(webSocketUrl, tickerName)
	defer wsConn.Close()
	aggService := service.AggregateService{}
	aggService.InitiateAggregateSequence(tickerName, aggregateTimeWindow, aggregatePersistentWindow, wsConn, nil)
}
