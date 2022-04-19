package application

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
	"time"
)

//TODO: Before going to production, the following would want to be things we might want to grab from env variables/volume
const webSocketUrl = "wss://delayed.polygon.io/stocks"

type Application struct {
	wsConn     service.WebSocketClient
	aggService *service.AggregateService
}

func NewApplication(tickerName string, level logrus.Level) *Application {
	logrus.SetLevel(level)
	webSocket := &service.TradeWebSocket{}
	webSocket.InitializeWSConnection(webSocketUrl, tickerName)
	application := &Application{
		wsConn:     webSocket,
		aggService: &service.AggregateService{},
	}

	return application
}

func (Application *Application) Run(tickerName string, aggregateTimeWindow time.Duration, aggregatePersistentWindow time.Duration) {
	Application.aggService.InitiateAggregateSequence(tickerName, aggregateTimeWindow, aggregatePersistentWindow, Application.wsConn, nil)
}
