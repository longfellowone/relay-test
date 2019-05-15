//go:generate go run github.com/99designs/gqlgen

// TODO: Move data loaders to new folder

// If user has scope manage:orders pull by user orginization id, else pull by user id

package dataloader

import (
	"context"
	"fmt"
	"time"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Project) IsNode() {}

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AddressID int
}

type Order struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Customer() CustomerResolver {
	return &customerResolver{r}
}

func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

type queryResolver struct{ *Resolver }

//SELECT * FROM customers
//SELECT * FROM address WHERE id IN (4,3,1)
//SELECT * FROM orders WHERE customer_id IN (2,1,3)

func (r *queryResolver) Customers(ctx context.Context) ([]*Customer, error) {
	fmt.Println("SELECT * FROM customers")
	//time.Sleep(5 * time.Millisecond)

	return []*Customer{
		{ID: 1, Name: "Bob", AddressID: 1},
		{ID: 2, Name: "Alice", AddressID: 3},
		{ID: 3, Name: "Eve", AddressID: 4},
	}, nil
}

func (r *queryResolver) Projects(ctx context.Context) ([]*Project, error) {
	fmt.Println("SELECT * FROM projects")

	time.Sleep(5 * time.Millisecond)

	return []*Project{
		{ID: 3, Name: "Project Name 3"},
		{ID: 5, Name: "Project Name 5"},
		{ID: 7, Name: "Project Name 7"},
		{ID: 2, Name: "Project Name 2"},
		{ID: 19, Name: "Project Name 19"},
	}, nil
}

func (r *queryResolver) Node(ctx context.Context, id string) (Node, error) {
	// Direct DB call
	return Project{ID: 99, Name: ""}, nil
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	return ctxLoaders(ctx).ordersByProject.Load(obj.ID)
}

type customerResolver struct{ *Resolver }

func (r *customerResolver) Orders(ctx context.Context, obj *Customer) ([]*Order, error) {
	return ctxLoaders(ctx).ordersByCustomer.Load(obj.ID)
}

//func (r *Resolver) Order() OrderResolver {
//	return &orderResolver{r}
//}

//func (r *customerResolver) Address(ctx context.Context, obj *Customer) (*Address, error) {
//	return ctxLoaders(ctx).addressByID.Load(obj.AddressID)
//}

//type orderResolver struct{ *Resolver }
//
//func (r *orderResolver) Items(ctx context.Context, obj *Order) ([]Item, error) {
//	return ctxLoaders(ctx).itemsByOrder.Load(obj.ID)
//}
