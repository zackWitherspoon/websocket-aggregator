package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"polygon-websocket-aggregator/application"
	"time"
)

//TODO: Before going to production, the following would want to be things we might want to grab from env variables/volume
const (
	aggregateTimeWindow       = 30 * time.Second
	aggregatePersistentWindow = 3600 * time.Second
)

func main() {
	//Get Ticker
	ticker := getTickerName()
	app := application.NewApplication(ticker, logrus.InfoLevel)
	app.Run(ticker, aggregateTimeWindow, aggregatePersistentWindow)
}

func getTickerName() string {
	if len(os.Args) != 2 {
		logrus.Error("Missing ticker name. Please include the ticker Name as the first argument.\n Example: ./main APPL")
		os.Exit(1)
	}
	return os.Args[1]
}
