package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
)

func Start(tickerName string, level logrus.Level) {
	logrus.SetLevel(level)
	webSockets := service.TradeWebSocket{}
	conn := webSockets.InitializeWSConnection(tickerName)
	defer conn.Close()

	service.InitiateAggregateCalculation(tickerName, conn)
}
