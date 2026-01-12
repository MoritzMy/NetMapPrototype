package main

import (
	"log"
	"net/http"

	"github.com/MoritzMy/NetMap/backend/internal/api"
	"github.com/MoritzMy/NetMap/backend/internal/graphing"
)

func main() {

	g := graphing.NewGraph()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/graph", api.GetGraph(g))
	http.HandleFunc("/api/icmp-sweep", api.RunICMPSweepHandler(g))
	http.HandleFunc("/api/arp-scan", api.RunARPScanHandler(g))

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
