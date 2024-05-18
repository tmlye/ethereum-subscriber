package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tmlye/ethereum-subscriber/pkg/blockpoller"
	"github.com/tmlye/ethereum-subscriber/pkg/ethgateway"
	"github.com/tmlye/ethereum-subscriber/pkg/storage"
)

func main() {
	const ethereumEndpoint = "https://cloudflare-eth.com"
	gateway := ethgateway.NewEthGateway(ethereumEndpoint)
	store := storage.NewMemoryStore()
	poller := blockpoller.NewBlockPoller(gateway, store)
	go poller.PollBlocks(12 * time.Second)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := strings.ToLower(r.URL.Query().Get("address"))
		if store.Subscribe(address) {
			fmt.Fprintf(w, "Subscribed to %s\n", address)
		} else {
			http.Error(w, "Already subscribed", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := strings.ToLower(r.URL.Query().Get("address"))
		transactions := store.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	})

	http.HandleFunc("/currentblock", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d\n", poller.LastProcessedBlock())
	})

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
