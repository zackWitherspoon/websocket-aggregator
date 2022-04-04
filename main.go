package main

import (
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/service"
)

const APIKEY = ""
const CHANNELS = "T.SPY,Q.SPY"

func main() {
	//Get Ticker
	ticker := "APPL"
	service.Start(ticker, logrus.DebugLevel)
}
