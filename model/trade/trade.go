package trade

import (
	"fmt"
)

type TradeRequest struct {
	Event      string  `json:"ev"`
	Symbol     string  `json:"sym"`
	ExchangeId int     `json:"x"`
	TradeId    string  `json:"i"`
	Tape       int     `json:"z"`
	Price      float64 `json:"p"`
	Size       int     `json:"s"`
	Conditions []int   `json:"c"`
	Timestamp  int64   `json:"t"`
}

func (trade *TradeRequest) PrintTrade() {
	conditionString := "["
	for i, v := range trade.Conditions {
		if i != len(trade.Conditions)-1 {
			conditionString += fmt.Sprintf("%d, ", v)
		} else {
			conditionString += fmt.Sprintf("%d ", v)
		}
	}
	conditionString += "]"

	fmt.Printf("{  \"ev\": \"%s\", \"sym\": \"%s\", \"x\": %d, \"i\": \"%s\", \"z\": %d, \"p\": %e, \"s\": %d, \"c\": %s, \"t\": %v }\n",
		trade.Event, trade.Symbol, trade.ExchangeId, trade.TradeId, trade.Tape, trade.Price, trade.Size, conditionString, trade.Timestamp)
}
