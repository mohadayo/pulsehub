package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mohadayo/pulsehub/services/gateway/config"
	"github.com/mohadayo/pulsehub/services/gateway/handler"
)

func main() {
	cfg := config.Load()
	h := handler.New(cfg.AnalyzerURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/events", h.ReceiveEvent)
	mux.HandleFunc("/stats", h.GetStats)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("[INFO] PulseHub Gateway starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
