//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OrderID   string `json:"order_ids"`
	Comments  string
	ProjectID string
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

type OrderConnection struct {
	IDs  []string
	From string
	To   string
}

func (c *OrderConnection) TotalCount() int {
	return len(c.IDs)
}

func (c *OrderConnection) PageInfo() PageInfo {
	return PageInfo{
		//StartCursor: EncodeCursor(c.From),
		//EndCursor:   EncodeCursor(c.To - 1),
		//HasNextPage: c.To < len(c.IDs),
	}
}

//func EncodeCursor(i string) string {
//	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1)))
//}

type Resolver struct {
	DB *sqlx.DB
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

func (r *Resolver) resolveOrderConnection(orderIDS []string, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	fmt.Println(orderIDS)

	from := 0
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i
	}

	to := len(orderIDS)
	if first != nil {
		to = from + *first
		if to > len(orderIDS) {
			to = len(orderIDS)
		}
	}

	//TODO: before, last

	fmt.Println(from, to)

	return &OrderConnection{}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Orders(ctx context.Context, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	// If user == purchaser load by org, else load ids by ctxUser
	//orderIDs := ctxLoaders(ctx).orderIdsByOrganization.Load(obj.ID)

	//return r.resolveOrderConnection(orderIDs, after, first, before, last)
	return &OrderConnection{}, nil
}

func (r *queryResolver) Projects(ctx context.Context) ([]*Project, error) {

	fmt.Println("SELECT * FROM projects WHERE orgID = ctxUserID")
	// Prime cache with project by id

	return []*Project{{ID: "projectID01"}, {ID: "projectID02"}}, nil
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	// Newest to oldest 120, 110, 100, 90, 80, 70

	orderIDs, err := ctxLoaders(ctx).orderIDsByProject.Load(obj.ID)
	if err != nil {
		fmt.Println(err)
	}

	return r.resolveOrderConnection(orderIDs, after, first, before, last)
}

// r.DB.Select(&orders, "SELECT * FROM orders WHERE sentdate < 110 order by sentdate DESC LIMIT 2")
