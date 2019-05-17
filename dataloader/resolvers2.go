//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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

type Resolver struct {
	DB *sqlx.DB
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

//func (r *queryResolver) Orders(ctx context.Context, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
//	// Newest to oldest 120, 110, 100, 90, 80, 70
//
//	// If the after argument is provided, add id > parsed_cursor to the WHERE clause
//	var orders []Order
//	err := r.DB.Select(&orders, "SELECT * FROM orders WHERE sentdate < 110 order by sentdate DESC LIMIT 2")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	for _, o := range orders {
//
//		fmt.Println(o.ID)
//	}
//
//	//return ctxLoaders(ctx).ordersByProject.Load(obj.ID)
//	return &OrderConnection{}, nil
//}

func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Projects(ctx context.Context) ([]*Project, error) {
	const query = `
		SELECT 
		id, name
		FROM projects
		LEFT JOIN orders on orders.projectid = id
-- 		WHERE projectid = $1
-- 		WHERE projectid IN ($1)
	`

	var projects []Project
	err := r.DB.Select(&projects, query)
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range projects {
		//fmt.Println(p.ID, p.Name, p.OrderID)
		fmt.Println(p)
	}

	return []*Project{{ID: "projectID01"}, {ID: "projectID02"}}, nil
}

//func (r *queryResolver) Projects(ctx context.Context, id string) (*ProjectConnection, error) {
//	const query = `
//		SELECT
//		orderid, projectid
//		FROM orders
//		LEFT JOIN projects on orders.projectid = id
//		WHERE projectid IN ($1,$2)
//	`
//
//	var projects []Project
//	err := r.DB.Select(&projects, query, id, "")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	for _, p := range projects {
//		//fmt.Println(p.ID, p.Name, p.OrderID)
//		fmt.Println(p.OrderID)
//	}
//
//	// No duplicate orders so join orders_id's on project
//
//	// Use load projects WHERE in orgID (From context)
//	//fmt.Println("SELECT * FROM project WHERE id = id")
//
//	// loader loaderOrderIds from project id
//	// Pass add order IDs to project struct
//
//	return &ProjectConnection{}, nil
//}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {

	fmt.Println(obj)

	//	const query = `
	//		SELECT
	//		id, name
	//		FROM projects
	//		LEFT JOIN orders on orders.projectid = id
	//		WHERE projectid = $1
	//-- 		WHERE projectid IN ($1)
	//	`
	//
	//	var projects []Project
	//	err := r.DB.Select(&projects, query, obj.ID)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	for _, p := range projects {
	//		//fmt.Println(p.ID, p.Name, p.OrderID)
	//		fmt.Println(p)
	//	}

	//fmt.Println(obj.OrderID)

	// Newest to oldest 120, 110, 100, 90, 80, 70

	// If the after argument is provided, add id > parsed_cursor to the WHERE clause
	//var orders []Order
	//err := r.DB.Select(&orders, "SELECT * FROM orders WHERE projectid=$1 AND sentdate < 110 order by sentdate DESC LIMIT 2", obj.ID)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//for _, o := range orders {
	//	fmt.Println(o.ID)
	//}

	// load orders by project ID
	//fmt.Println("Select * FROM orders WHERE projectID in (2,8,12)")

	//return ctxLoaders(ctx).ordersByProject.Load(obj.ID)
	return &OrderConnection{}, nil
}

// Prime project to order loader
//const resolvers = {
//Query: {
//users: async (root, args, { userLoader }) => {
//// Our loader can't get all users, so let's use the model directly here
//const allUsers = await User.find({})
//// then tell the loader about the users we found
//for (const user of allUsers) {
//userLoader.prime(user.id, user);
//}
//// and finally return the result
//return allUsers
//}
//},
//User: {
//friends: async (user, args, { userLoader }) => {
//return userLoader.loadMany(user.friendIds)
//},
//},
//}
