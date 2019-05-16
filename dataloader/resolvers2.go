//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"fmt"
	"time"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Project) IsNode() {}

type Order struct {
	ID     string    `json:"id"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

func (Order) IsNode() {}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Projects(ctx context.Context) ([]*Project, error) {
	fmt.Println("SELECT * FROM projects")

	time.Sleep(5 * time.Millisecond)

	return []*Project{
		{ID: "3", Name: "Project Name 3"},
		{ID: "5", Name: "Project Name 5"},
		{ID: "7", Name: "Project Name 7"},
		{ID: "2", Name: "Project Name 2"},
		{ID: "19", Name: "Project Name 19"},
	}, nil
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	return ctxLoaders(ctx).ordersByProject.Load(obj.ID)
}
