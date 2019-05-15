//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"fmt"
	"time"
)

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

func (r *Resolver) Customer() CustomerResolver {
	return &customerResolver{r}
}

func (r *Resolver) Order() OrderResolver {
	return &orderResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type customerResolver struct{ *Resolver }

func (r *customerResolver) Address(ctx context.Context, obj *Customer) (*Address, error) {
	return ctxLoaders(ctx).addressByID.Load(obj.AddressID)
}

func (r *customerResolver) Orders(ctx context.Context, obj *Customer) ([]Order, error) {
	return ctxLoaders(ctx).ordersByCustomer.Load(obj.ID)
}

type orderResolver struct{ *Resolver }

func (r *orderResolver) Items(ctx context.Context, obj *Order) ([]Item, error) {
	return ctxLoaders(ctx).itemsByOrder.Load(obj.ID)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Customers(ctx context.Context) ([]Customer, error) {
	fmt.Println("SELECT * FROM customer")

	time.Sleep(5 * time.Millisecond)

	return []Customer{
		{ID: 1, Name: "Bob", AddressID: 1},
		{ID: 2, Name: "Alice", AddressID: 3},
		{ID: 3, Name: "Eve", AddressID: 4},
	}, nil
}
