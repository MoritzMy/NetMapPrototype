package _map

type EdgeType string

const (
	EdgeMemberOf   EdgeType = "member-of"
	EdgeRouteVia   EdgeType = "routes-via"
	EdgeRespondsTo EdgeType = "responds-to"
)

type Edge struct {
	From *Node
	To   *Node
	Type EdgeType
}
