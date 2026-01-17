package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/MoritzMy/NetMap/backend/internal/api"
	"github.com/MoritzMy/NetMap/backend/internal/graphing"
)

const (
	defaultPort       = 8080
	maxPortsInNetwork = 65535
)

func main() {

	port := flag.Uint("p", defaultPort, "Sets the Port for the HTTP Server to listen on")

	flag.Parse()

	if *port > maxPortsInNetwork {
		*port = defaultPort // reset to default value
	}

	addr := ":" + strconv.Itoa(int(*port))

	g := graphing.NewGraph()
	graphing.CreateLocalHostNetworkNodes(g)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/graph", api.GetGraph(g))
	http.HandleFunc("/api/icmp-sweep", api.RunICMPSweepHandler(g))
	http.HandleFunc("/api/arp-scan", api.RunARPScanHandler(g))
	http.HandleFunc("/api/reset", api.ResetGraph(g))

	log.Println("Listening on port 8080")
	log.Fatal("service crashed with: ", http.ListenAndServe(addr, nil))
}
