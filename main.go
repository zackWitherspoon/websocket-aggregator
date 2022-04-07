package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"polygon-websocket-aggregator/application"
)

const APIKEY = ""

func main() {
	//Get Ticker
	if len(os.Args) != 2 {
		logrus.Error("Missing ticker name. Please include the ticker Name as the first argument.\n Example: ./main APPL")
		os.Exit(1)
	}
	ticker := os.Args[1]
	application.Start(ticker, logrus.InfoLevel)
}
