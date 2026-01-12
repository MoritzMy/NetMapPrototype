package api

import (
	"net/http"

	"github.com/MoritzMy/NetMap/backend/internal/arp_scan"
	"github.com/MoritzMy/NetMap/backend/internal/graphing"
	"github.com/MoritzMy/NetMap/backend/internal/ping"
)

func RunICMPSweepHandler(g *graphing.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Start ICMP Sweep

		go func() {
			ping.RunICMPSweep(g)
		}()

		w.Write([]byte("Started ICMP Sweep"))
	}
}

func RunARPScanHandler(g *graphing.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Start ARP Scan

		go func() {
			arp_scan.RunARPScan(g)
		}()

		w.Write([]byte("Started ARP Scan"))
	}
}
