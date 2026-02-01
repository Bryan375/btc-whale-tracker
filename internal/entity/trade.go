package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		Symbol    string          `json:"symbol"`
		Price     decimal.Decimal `json:"price"`
		Quantity  decimal.Decimal `json:"quantity"`
		Timestamp time.Time       `json:"timestamp"`
	}

	WhaleAlert struct {
		ID            int64           `json:"id"`
		Symbol        string          `json:"symbol"`
		AmountUSD     decimal.Decimal `json:"amount_usd"`
		Action        string          `json:"action"`
		AIExplanation string          `json:"ai_explanation"`
		Timestamp     time.Time       `json:"timestamp"`
	}
)
