package server

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) AddTodo(ctx context.Context, input AddTodoInput) (*AddTodoPayload, error) {
	panic("not implemented")
}
func (r *mutationResolver) ChangeTodoStatus(ctx context.Context, input ChangeTodoStatusInput) (*ChangeTodoStatusPayload, error) {
	panic("not implemented")
}
func (r *mutationResolver) MarkAllTodos(ctx context.Context, input MarkAllTodosInput) (*MarkAllTodosPayload, error) {
	panic("not implemented")
}
func (r *mutationResolver) RemoveCompletedTodos(ctx context.Context, input RemoveCompletedTodosInput) (*RemoveCompletedTodosPayload, error) {
	panic("not implemented")
}
func (r *mutationResolver) RemoveTodo(ctx context.Context, input RemoveTodoInput) (*RemoveTodoPayload, error) {
	panic("not implemented")
}
func (r *mutationResolver) RenameTodo(ctx context.Context, input RenameTodoInput) (*RenameTodoPayload, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string) (*User, error) {
	panic("not implemented")
}
func (r *queryResolver) Node(ctx context.Context, id string) (Node, error) {
	panic("not implemented")
}
