package _map

import "sync"

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

func (g *Graph) GetOrCreateNode(id string, nodeType NodeType) *Node {
	g.mu.Lock()
	defer g.mu.Unlock()

	if node, exists := g.Nodes[id]; exists {
		return &node
	}

	newNode := Node{
		ID:   id,
		Type: nodeType,
	}
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
