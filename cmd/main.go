package main

import (
	"log"

	"github.com/Bryan375/btc-whale-tracker/internal/ingest"
)

func main() {
	client := ingest.NewTokocryptoClient()
	if err := client.Connect(); err != nil {
		log.Fatal("Connection failed:", err)
	}

	svc := ingest.NewService(client, 10000.0)
	svc.Start()
}
