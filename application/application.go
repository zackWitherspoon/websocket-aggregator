package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
)

const webSocketUrl = "wss://delayed.polygon.io/stocks"

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSockets := service.TradeWebSocket{}
	conn := webSockets.InitializeWSConnection(webSocketUrl, tickerName)
	defer conn.Close()
	aggService := service.AggregateService{}
	aggService.InitiateAggregateCalculation(tickerName, conn, nil)
}
