package models

import (
	"time"
)

// PurchaseRequest represents the payload for creating a new crypto purchase
// Note: RemainingQty will be set equal to Quantity on creation (nothing sold yet)
type PurchaseRequest struct {
	Symbol       string     `json:"symbol" binding:"required"`                   // e.g. 'ETHUSDT' or 'ETH'
	Exchange     string     `json:"exchange" binding:"required"`                 // e.g. 'binance'
	OrderID      string     `json:"order_id" binding:"required"`                 // exchange order id for traceability
	Quantity     float64    `json:"quantity" binding:"required,gt=0"`            // crypto units purchased
	FiatInvested float64    `json:"fiat_invested" binding:"required,gt=0"`       // USD amount invested
	BuyPrice     float64    `json:"buy_price" binding:"required,gt=0"`           // price per unit at buy
	Fee          float64    `json:"fee,omitempty"`                               // transaction fee
	FeeCurrency  string     `json:"fee_currency,omitempty"`                      // currency of fee
	TargetGain   float64    `json:"target_gain" binding:"omitempty,gte=0,lte=1"` // default 0.03 (3%)
	Status       string     `json:"status,omitempty"`                            // 'open', 'partial', 'closed'
	CreatedAt    *time.Time `json:"created_at,omitempty"`                        // optional timestamp
}

// SellRequest represents the payload for recording a crypto sell
type SellRequest struct {
	PurchaseID  int64      `json:"purchase_id" binding:"required"`      // references crypto_purchases.id
	Exchange    string     `json:"exchange" binding:"required"`         // e.g. 'binance'
	OrderID     string     `json:"order_id" binding:"required"`         // exchange order id for traceability
	Quantity    float64    `json:"quantity" binding:"required,gt=0"`    // crypto units sold
	SellPrice   float64    `json:"sell_price" binding:"required,gt=0"`  // price per unit at sell
	FiatAmount  float64    `json:"fiat_amount" binding:"required,gt=0"` // total fiat received
	Fee         float64    `json:"fee,omitempty"`                       // transaction fee
	FeeCurrency string     `json:"fee_currency,omitempty"`              // currency of fee
	ExecutedAt  *time.Time `json:"executed_at,omitempty"`               // optional timestamp
}
