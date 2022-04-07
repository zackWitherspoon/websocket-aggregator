package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
)

const (
	webSocketUrl        = "wss://delayed.polygon.io/stocks"
	aggregateTime       = 5 * SecondAlteration
	aggregateKeepWindow = 3600 * SecondAlteration
	SecondAlteration    = 1000000000
)

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSockets := service.TradeWebSocket{}
	conn := webSockets.InitializeWSConnection(webSocketUrl, tickerName)
	defer conn.Close()
	aggService := service.AggregateService{}
	aggService.InitiateAggregateCalculation(tickerName, aggregateTime, aggregateKeepWindow, conn, nil)
}
