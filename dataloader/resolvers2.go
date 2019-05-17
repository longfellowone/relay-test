//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Project) IsNode() {}

type Order struct {
	ID           string `db:"orderid"`
	ProjectId    string `db:"projectid"`
	Status       int    `db:"status"`
	ProjectName  string `db:"project_name"`
	ForemanEmail string `db:"foreman_email"`
	SentDate     int64  `db:"sentdate"`
	Comments     string `db:"comments"`
}

func (Order) IsNode() {}

type Resolver struct {
	DB *sqlx.DB
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Project(ctx context.Context, id string) (*Project, error) {
	fmt.Println("SELECT * FROM project WHERE id = id")

	return &Project{ID: "cf510766-faf7-415e-a067-0c5ae5cb2ae8", Name: "Project Name 3"}, nil
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	// Fetch first 10 timestamps for a project
	// pass array of timestamps in connection to edges
	// Edges resolveOrdersByProject, pass in array of timestamps
	// SELECT * FROM orders where timestamp IN -----
	// array of project order IDs

	var orders []Order
	err := r.DB.Select(&orders, "SELECT sentdate FROM orders WHERE projectid=$1 order by sentdate ASC", obj.ID)
	if err != nil {
		fmt.Println(err)
	}

	for _, o := range orders {
		fmt.Println(o.SentDate)
	}

	// 	fmt.Println("SELECT timestamp FROM orders WHERE project_id = obj.id LIMIT first")
	//return &OrderConnection{
	//	PageInfo: &PageInfo{
	//		HasNextPage:     false,
	//		HasPreviousPage: false,
	//		StartCursor:     nil,
	//		EndCursor:       nil,
	//	},
	//	Edges: nil,
	//}, nil
	//return ctxLoaders(ctx).ordersByProject.Load(obj.ID)
	// Dont even use data loader?
	return &OrderConnection{}, nil
}
