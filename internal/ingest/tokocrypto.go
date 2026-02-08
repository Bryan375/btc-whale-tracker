package ingest

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Bryan375/btc-whale-tracker/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const (
	TokocryptoWSURL = "wss://stream-cloud.tokocrypto.site/stream"

	MethodSubscribe = "SUBSCRIBE"

	EventTypeAggTrade = "aggTrade"

	SymbolBTCUSDT = "BTCUSDT"
)

type (
	TokocryptoSubscribeMsg struct {
		Method string   `json:"method"`
		Params []string `json:"params"`
		ID     int      `json:"id"`
	}

	TokocryptoResponse struct {
		Stream string                  `json:"stream"`
		Data   TokocryptoAggTradeEvent `json:"data"`
	}

	TokocryptoAggTradeEvent struct {
		EventType    string `json:"e"`
		EventTime    int64  `json:"E"`
		Symbol       string `json:"s"`
		AggregateID  int64  `json:"a"`
		Price        string `json:"p"`
		Quantity     string `json:"q"`
		FirstTradeID int64  `json:"f"`
		LastTradeID  int64  `json:"l"`
		Timestamp    int64  `json:"T"`
		IsBuyer      bool   `json:"m"`
		IsBestMatch  bool   `json:"M"`
	}

	TokocryptoClient struct {
		conn *websocket.Conn
		url  string
	}
)

func NewTokocryptoClient() *TokocryptoClient {
	return &TokocryptoClient{
		url: TokocryptoWSURL,
	}
}

func (c *TokocryptoClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("tokocrypto dial failed: %w", err)
	}
	c.conn = conn

	subscribeMsg := TokocryptoSubscribeMsg{
		Method: MethodSubscribe,
		Params: []string{"btcusdt@aggTrade"},
		ID:     1,
	}

	if err := c.conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("tokocrypto subscribe failed: %w", err)
	}

	return nil
}

func (c *TokocryptoClient) closeConnection() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Tokocrypto connection close error: %v", err)
	}
}

func (c *TokocryptoClient) Stream(tradeChan chan<- entity.Trade) {
	defer c.closeConnection()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		log.Printf("📨 Raw message: %s", string(message))

		var response TokocryptoResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("JSON error: %v", err)
			continue
		}

		event := response.Data
		if event.EventType != EventTypeAggTrade {
			continue
		}

		price, _ := decimal.NewFromString(event.Price)
		qty, _ := decimal.NewFromString(event.Quantity)
		timestamp := time.Unix(0, event.Timestamp*int64(time.Millisecond))

		trade := entity.Trade{
			Symbol:    SymbolBTCUSDT,
			Price:     price,
			Quantity:  qty,
			Timestamp: timestamp,
		}

		tradeChan <- trade
	}
}
