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
}
