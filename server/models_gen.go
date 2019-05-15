// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package server

type Node interface {
	IsNode()
}

type Order struct {
	ID      string   `json:"id"`
	Project Project  `json:"project"`
	Items   []string `json:"items"`
}

func (Order) IsNode() {}

type OrderConnection struct {
	PageInfo PageInfo     `json:"pageInfo"`
	Edges    []*OrderEdge `json:"edges"`
}

type OrderEdge struct {
	Node   *Order `json:"node"`
	Cursor string `json:"cursor"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Project) IsNode() {}
