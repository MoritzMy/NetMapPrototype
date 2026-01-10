package graphing

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Graph struct {
	Edges []*Edge
	Nodes map[string]*Node
	mu    sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{
		Edges: make([]*Edge, 0),
		Nodes: make(map[string]*Node),
	}
}

func (g *Graph) GetOrCreateNode(id string) *Node {
	g.mu.Lock()
	defer g.mu.Unlock()

	if node, exists := g.Nodes[id]; exists {
		return node
	}

	newNode := newNode(id)

	g.Nodes[id] = newNode
	return newNode
}

func (g *Graph) AddEdge(fromID, toID string, edgeType EdgeType) {
	g.mu.Lock()
	defer g.mu.Unlock()

	fromNode, fromExists := g.Nodes[fromID]
	toNode, toExists := g.Nodes[toID]

	if !fromExists || !toExists {
		return
	}

	edge := &Edge{
		From: fromNode,
		To:   toNode,
		Type: edgeType,
	}
	g.Edges = append(g.Edges, edge)
}

func (g *Graph) AddProtocol(id, proto string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if node, exists := g.Nodes[id]; exists {
		node.Protocols[proto] = true
	}
}

func (g *Graph) LinkNetworkToGateway() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, edge := range g.Edges {
		if edge.Type == EdgeMemberOf {
			if edge.From.Type == NodeNetwork && edge.To.Type == NodeGateway {
				routeEdge := &Edge{
					From: edge.From,
					To:   edge.To,
					Type: EdgeRouteVia,
				}
				g.Edges = append(g.Edges, routeEdge)
			} else if edge.To.Type == NodeNetwork && edge.From.Type == NodeGateway {
				routeEdge := &Edge{
					From: edge.To,
					To:   edge.From,
					Type: EdgeRouteVia,
				}
				g.Edges = append(g.Edges, routeEdge)
			}
		}
	}
}

func (g *Graph) String() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	result := "Graph:\n"
	result += "Nodes:\n"
	for _, node := range g.Nodes {
		result += fmt.Sprintf("- ID: %s responded to: %v and has Type: %v\n", node.ID, node.Protocols, node.Type)
	}
	result += "Edges:\n"
	for _, edge := range g.Edges {
		result += "- From: " + edge.From.ID + " To: " + edge.To.ID + " Type: " + string(edge.Type) + "\n"
	}

	return result
}

func (g *Graph) ToDOT() string {
	var b strings.Builder

	b.WriteString("digraph netmap {\n")
	b.WriteString("  rankdir=LR;\n")
	b.WriteString("  node [fontname=\"Helvetica\"];\n\n")

	// Nodes
	for _, n := range g.Nodes {
		shape := "ellipse"
		color := "black"

		switch n.Type {
		case NodeHost:
			shape = "box"
		case NodeNetwork:
			shape = "oval"
		case NodeGateway:
			shape = "diamond"
		}

		label := n.ID
		if n.IP != nil {
			label = n.IP.String()
		}

		fmt.Fprintf(
			&b,
			"  \"%s\" [label=\"%s\", shape=%s, color=%s];\n",
			n.ID, label, shape, color,
		)
	}

	b.WriteString("\n")

	// Edges
	for _, e := range g.Edges {
		fmt.Fprintf(
			&b,
			"  \"%s\" -> \"%s\" [label=\"%s\"];\n",
			e.From.ID, e.To.ID, e.Type,
		)
	}

	b.WriteString("}\n")
	return b.String()
}

func (g *Graph) ExportToDOT(filename string) error {
	dot := g.ToDOT()
	err := os.WriteFile(filename+".dot", []byte(dot), 0644)
	return err
}

type graphJSON struct {
	Nodes []graphNode `json:"nodes"`
	Links []graphLink `json:"links"`
}

type graphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

type graphNode struct {
	ID         string          `json:"id"`
	Type       NodeType        `json:"type"`
	IP         net.IP          `json:"ip"`
	MAC        string          `json:"mac"`
	Protocols  map[string]bool `json:"protocols"`
	Vendor     string          `json:"vendor"`
	Confidence float64         `json:"confidence"`
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	nodes := make([]graphNode, 0, len(g.Nodes))

	for _, n := range g.Nodes {
		nodes = append(nodes, graphNode{
			ID:         n.ID,
			Type:       n.Type,
			IP:         n.IP,
			MAC:        n.MAC.String(),
			Protocols:  n.Protocols,
			Vendor:     n.Vendor,
			Confidence: n.Confidence,
		})
	}

	links := make([]graphLink, 0, len(g.Edges))
	for _, e := range g.Edges {
		if e.From == nil || e.To == nil {
			return nil, fmt.Errorf("edge cannot be marshalled: missing from or to")
		}

		links = append(links, graphLink{
			Source: e.From.ID,
			Target: e.To.ID,
			Type:   string(e.Type),
		})
	}

	return json.Marshal(graphJSON{
		Nodes: nodes,
		Links: links,
	})
}
