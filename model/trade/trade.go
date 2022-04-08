package trade

type TradeRequest struct {
	Symbol    string  `json:"sym"`
	Price     float64 `json:"p"`
	Size      int     `json:"s"`
	Timestamp int64   `json:"t"`
}
