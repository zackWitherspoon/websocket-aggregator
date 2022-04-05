package trade

import (
	"encoding/json"
	"fmt"
)

type TradeRequest struct {
	Symbol    string  `json:"sym"`
	TradeId   string  `json:"i"`
	Price     float64 `json:"p"`
	Size      int     `json:"s"`
	Timestamp int64   `json:"t"`
}

func (trade *TradeRequest) PrintTrade() {
	res, err := json.Marshal(trade)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
	//fmt.Printf("{  \"sym\": \"%s\", \"i\": \"%s\", \"z\": %d, \"p\": %e, \"s\": %d, \"c\": %s, \"t\": %v }\n",
	//	v.Symbol, trade.ExchangeId, trade.TradeId, trade.Tape, trade.Price, trade.Size, conditionString, trade.Timestamp)
}
