package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
	"time"
)

const (
	webSocketUrl        = "wss://delayed.polygon.io/stocks"
	aggregateTime       = 5 * secondAlteration
	aggregateKeepWindow = 3600 * secondAlteration
	secondAlteration    = 1000000000
)

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSockets := service.TradeWebSocket{}
	conn := webSockets.InitializeWSConnection(webSocketUrl, tickerName)
	defer conn.Close()
	aggService := service.AggregateService{}
	println("aggregateTime = ", aggregateTime)
	println("aggregateTime in seconds = ", aggregateTime*time.Second)
	aggService.InitiateAggregateCalculation(tickerName, aggregateTime, aggregateKeepWindow, conn, nil)
}
