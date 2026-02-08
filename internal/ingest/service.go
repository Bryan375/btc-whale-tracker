package ingest

import (
	"log"

	"github.com/Bryan375/btc-whale-tracker/internal/entity"
	"github.com/shopspring/decimal"
)

type Service struct {
	client    *TokocryptoClient
	tradeChan chan entity.Trade
	minValue  decimal.Decimal
}

func NewService(client *TokocryptoClient, minUSD float64) *Service {
	return &Service{
		client:    client,
		tradeChan: make(chan entity.Trade, 100),
		minValue:  decimal.NewFromFloat(minUSD),
	}
}

func (s *Service) Start() {
	go s.client.Stream(s.tradeChan)

	log.Printf("Whale Tracker Started (Filter: >$%s)", s.minValue.String())

	for trade := range s.tradeChan {
		log.Printf("Received Trade: %+v", trade)
		value := trade.Price.Mul(trade.Quantity)

		if value.GreaterThanOrEqual(s.minValue) {
			log.Printf("🐋 WHALE ALERT: %s | Qty: %s | Val: $%s",
				trade.Symbol,
				trade.Quantity.StringFixed(4),
				value.StringFixed(2))

			// TODO: Publish to RabbitMQ (FR-03)
		}
	}
}
