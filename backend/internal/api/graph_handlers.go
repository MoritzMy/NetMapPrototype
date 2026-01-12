package api

import (
	"net/http"

	"github.com/MoritzMy/NetMap/backend/internal/graphing"
)

// GetGraph is an HTTP handler that returns the graph as JSON.
func GetGraph(g *graphing.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json, err := g.MarshalJSON() // Serialize the graph to JSON
		if err != nil {
			http.Error(w, "Could not create json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(json)
		if err != nil {
			return
		}
	}
}

func ResetGraph(g *graphing.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		g = graphing.NewGraph()

		w.Write([]byte("Graph reset"))
	}
}
