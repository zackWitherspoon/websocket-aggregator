package main

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/application"
)

const APIKEY = ""
const CHANNELS = "T.SPY,Q.SPY"

func main() {
	//Get Ticker
	ticker := "APPL"
	application.Start(ticker, logrus.DebugLevel)
}
