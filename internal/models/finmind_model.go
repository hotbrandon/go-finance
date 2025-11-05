package models

import (
	"encoding/json"
)

type FinmindResponse struct {
	Msg    string          `json:"message"`
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
}

// typed model for the TaiwanStockPrice dataset
type StockPrice struct {
	Date            string  `json:"date"`
	StockID         string  `json:"stock_id"`
	TradingVolume   int64   `json:"Trading_Volume"`
	TradingMoney    int64   `json:"Trading_money"`
	Open            float64 `json:"open"`
	Max             float64 `json:"max"`
	Min             float64 `json:"min"`
	Close           float64 `json:"close"`
	Spread          float64 `json:"spread"`
	TradingTurnover int64   `json:"Trading_turnover"`
}

// typed response for this specific endpoint (use in docs/handlers)
type TaiwanStockPriceResponse struct {
	Msg    string       `json:"message"`
	Status int          `json:"status"`
	Data   []StockPrice `json:"data"`
}
