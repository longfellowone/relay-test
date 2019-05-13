package server

import (
	"context"
	"fmt"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) TodoConnection() TodoConnectionResolver {
	return &todoConnectionResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (r *Resolver) resolveTodos(ids []string) ([]*Todo, error) {
	fmt.Println("r/resolvetodo")
	return []*Todo{{
		ID:       "this test id for a todo",
		Text:     "",
		Complete: false,
	}}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string) (*User, error) {
	fmt.Println("query/user")
	return &User{
		ID: "test user id",
		TotalCount:     0,
		CompletedCount: 0,
	}, nil
}
func (r *queryResolver) Node(ctx context.Context, id string) (Node, error) {
	fmt.Println("get type from global id: ", id)
	return &Todo{
		ID:       "a single todo",
		Text:     "",
		Complete: false,
	}, nil
}

type todoConnectionResolver struct{ *Resolver }

func (r *todoConnectionResolver) Edges(ctx context.Context, obj *TodoConnection) ([]*TodoEdge, error) {
	todos, err := r.resolveTodos([]string{"1","2"})
	if err != nil {
		return []*TodoEdge{}, nil
	}

	fmt.Println("todoconn/edges", obj.Edges)

	edges := make([]*TodoEdge, 2)
	for i := range edges {
		edges[i] = &TodoEdge{
			Node:   todos[0],
			Cursor: "my cursor",
		}

	}

	return edges,nil
}


func (r *todoConnectionResolver) PageInfo(ctx context.Context, obj *TodoConnection) (*PageInfo, error) {
	return &PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor:     nil,
		EndCursor:       nil,
	}, nil
}


type userResolver struct{ *Resolver }

func (r *userResolver) Todos(ctx context.Context, obj *User, after *string, first *int, before *string, last *int) (*TodoConnection, error) {
	// obj.IdArray
	fmt.Println("user/todos")
	return &TodoConnection{}, nil
}
