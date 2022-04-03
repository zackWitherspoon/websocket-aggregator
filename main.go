package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const APIKEY = "JncA2_tG82oqnPqGQHyMj2BW_h2jCoQl"
const CHANNELS = "T.SPY,Q.SPY"

func main() {

	//connect to websocket
	c, _, err := websocket.DefaultDialer.Dial("wss://delayed.polygon.io/stocks", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	//authenticate websocket
	auth := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKEY)))

	//subscribe to websocket
	subscribe := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", CHANNELS)))
	if auth != nil {
		println("ERROR", auth)
	} else {
		println("NO ERROR. We are authenticated")
	}
	if subscribe != nil {
		println("ERROR", subscribe)
	} else {
		println("NO ERROR. We are Subscribed!")
	}

	// Buffered channel to account for bursts or spikes in data:
	chanMessages := make(chan interface{}, 10000)

	// Read messages off the buffered queue:
	go func() {
		for msgBytes := range chanMessages {
			logrus.Info("Message Bytes: ", msgBytes)
		}
	}()

	// As little logic as possible in the reader loop:
	for {
		var msg interface{}
		//print websocket
		err := c.ReadJSON(&msg)
		// Ideally use c.ReadMessage instead of ReadJSON so you can parse the JSON data in trade
		// separate go routine. Any processing done in this loop increases the chances of disconnects
		// due to not consuming the data fast enough.
		if err != nil {
			panic(err)
		}
		chanMessages <- msg
	}
}
