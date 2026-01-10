package graphing

type EdgeType string

const (
	EdgeMemberOf   EdgeType = "member-of"
	EdgeRouteVia   EdgeType = "routes-via"
	EdgeRespondsTo EdgeType = "responds-to"
)

type Edge struct {
	From *Node    `json:"-"`
	To   *Node    `json:"-"`
	Type EdgeType `json:"type"`
}
