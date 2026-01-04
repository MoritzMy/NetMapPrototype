package _map

import (
	"fmt"
	"sync"
)

type Graph struct {
	Edges []*Edge
	Nodes map[string]Node
	mu    sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{
		Edges: make([]*Edge, 0),
		Nodes: make(map[string]Node),
	}
}

func (g *Graph) GetOrCreateNode(id string) *Node {
	g.mu.Lock()
	defer g.mu.Unlock()

	if node, exists := g.Nodes[id]; exists {
		return &node
	}

	newNode := newNode(id)

	g.Nodes[id] = newNode
	return &newNode
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
		From: &fromNode,
		To:   &toNode,
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

func (g *Graph) String() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	result := "Graph:\n"
	result += "Nodes:\n"
	for _, node := range g.Nodes {
		result += fmt.Sprintf("- ID: %s responded to: %v\n", node.ID, node.Protocols)
	}
	result += "Edges:\n"
	for _, edge := range g.Edges {
		result += "- From: " + edge.From.ID + " To: " + edge.To.ID + " Type: " + string(edge.Type) + "\n"
	}
	return result
}
