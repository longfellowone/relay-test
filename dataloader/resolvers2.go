//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

//  return errors.Wrap(err, "request failed")

type Project struct {
	ID       string
	Name     string
	OrderIDs []string
	//Comments  string
	//ProjectID string
}

func (Project) IsNode() {}

//SELECT `books`.title,
//`book_order`.orderdate,
//`book_order`.orderid
//FROM   `books`
//LEFT OUTER JOIN `order_items`
//ON `books`.bookid = `order_items`.bookid
//LEFT OUTER JOIN `book_order`
//ON `order_items`.orderid = `book_order`.orderid
//ORDER  BY `order_items`.bookid DESC;

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
	TotalCount int `json:"totalCount"`
	PageInfo   *PageInfo
	IDs        []string
	Begin      int
	End        int
}

//func (c *OrderConnection) PageInfo() PageInfo {
//
//	return PageInfo{
//		//StartCursor: EncodeCursor(c.From),
//		//EndCursor:   EncodeCursor(c.To - 1),
//		//HasNextPage: c.To < len(c.IDs),
//	}
//}

//func EncodeCursor(i string) string {
//	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1)))
//}

type Resolver struct {
	DB *sqlx.DB
}

func (r *Resolver) OrderConnection() OrderConnectionResolver {
	return &orderConnectionResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Project() ProjectResolver {
	return &projectResolver{r}
}

type orderConnectionResolver struct{ *Resolver }

func (r *orderConnectionResolver) Edges(ctx context.Context, obj *OrderConnection) ([]*OrderEdge, error) {
	nodes, err := r.resolveOrders(ctx, obj.IDs)
	if err != nil {
		fmt.Println(err)
	}

	edges := make([]*OrderEdge, obj.End-obj.Begin)

	for i := range edges {
		edges[i] = &OrderEdge{
			Node:   nodes[i],
			Cursor: strconv.Itoa(obj.Begin + i),
		}
	}

	return edges, nil
}

func (r *Resolver) resolveOrders(ctx context.Context, ids []string) ([]*Order, error) {
	orders, errs := ctxLoaders(ctx).orderLoader.LoadAll(ids)
	for _, err := range errs {
		if err != nil {
			return []*Order{}, err
		}
	}

	return orders, nil
}

type Filter map[string]interface{}

//resolve: (user, { first }) => queryLoader.load([
//'SELECT toID FROM friends WHERE fromID=? LIMIT ?', user.id, first
//]).then(rows => rows.map(row => userLoader.load(row.toID)))

func (r *Resolver) resolveOrderConnection(orderIDs []string, after *string, first *int, before *string, last *int) (*OrderConnection, error) {

	// loader
	// edgeBuilder

	temp := func(filter Filter) ([]string, int, int, *PageInfo) {
		args := NewConnectionArguments(filter)

		arraySlice := orderIDs
		arrayLength := len(arraySlice)

		beforeOffset := GetOffsetWithDefault(args.Before, arrayLength)
		afterOffset := GetOffsetWithDefault(args.After, -1)

		startOffset := max(afterOffset, -1) + 1
		endOffset := min(beforeOffset, arrayLength)

		if args.First != -1 {
			endOffset = min(endOffset, startOffset+args.First)
		}

		if args.Last != -1 {
			startOffset = max(startOffset, endOffset-args.Last)
		}

		begin := startOffset
		end := arrayLength - (arrayLength - endOffset)

		if begin > end {
			return nil, 0, 0, &PageInfo{}
		}

		slice := arraySlice[begin:end]

		//var firstEdgeCursor, lastEdgeCursor ConnectionCursor
		var firstEdgeCursor, lastEdgeCursor *string
		if len(slice) > 0 {
			firstEdgeCursor = OffsetToCursor(startOffset)
			lastEdgeCursor = OffsetToCursor(startOffset + len(slice) - 1)
		}

		hasPreviousPage := false
		if startOffset > 0 {
			hasPreviousPage = true
		}

		hasNextPage := false
		if endOffset < arrayLength {
			hasNextPage = true
		}

		pageInfo := &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstEdgeCursor,
			EndCursor:       lastEdgeCursor,
		}

		return slice, begin, end, pageInfo
	}

	filter := map[string]interface{}{
		"first": 4,
		"after": "cursor:0",
	}

	ids, begin, end, pageInfo := temp(filter)

	//count := len(orderIDs)

	return &OrderConnection{
		TotalCount: len(orderIDs),
		PageInfo:   pageInfo,
		IDs:        ids,
		Begin:      begin,
		End:        end,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Orders(ctx context.Context, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	// If user == purchaser load by org, else load ids by ctxUser
	//orderIDs := ctxLoaders(ctx).orderIdsByOrganization.Load(obj.ID)

	// DO NO USE SELECT WITHOUT LIMIT
	var ids []string
	err := r.DB.Select(&ids, "SELECT orderid FROM orders ORDER BY sentdate DESC")
	if err != nil {
		return &OrderConnection{}, err
	}

	return r.resolveOrderConnection(ids, after, first, before, last)
}

func (r *queryResolver) Projects(ctx context.Context) ([]*Project, error) {

	fmt.Println("SELECT * FROM projects WHERE orgID = ctxUserOrgID")
	// Prime cache with project by id

	const query = `
		SELECT p.id, p.name, array_agg(o.orderid ORDER BY o.sentdate DESC) as order_ids
		FROM projects p INNER JOIN orders o on p.id = o.projectid
		GROUP BY p.id
		ORDER BY p.name ASC 
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	var id string
	var name string
	var order_ids []string

	var projects []*Project

	for rows.Next() {
		if err := rows.Scan(&id, &name, pq.Array(&order_ids)); err != nil {
			fmt.Println(err)
		}

		projects = append(projects, &Project{
			ID:       id,
			Name:     name,
			OrderIDs: order_ids,
		})
	}

	return projects, nil
}

type projectResolver struct{ *Resolver }

func (r *projectResolver) Orders(ctx context.Context, obj *Project, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
	// Newest to oldest 120, 110, 100, 90, 80, 70

	//orderIDs, err := ctxLoaders(ctx).orderIDsByProject.Load(obj.ID)
	//if err != nil {
	//	fmt.Println(err)
	//}

	return r.resolveOrderConnection(obj.OrderIDs, after, first, before, last)
}

const PREFIX = "cursor:"

type ConnectionCursor string

type ConnectionArguments struct {
	Before ConnectionCursor
	After  ConnectionCursor
	First  int
	Last   int
}

func NewConnectionArguments(filters map[string]interface{}) ConnectionArguments {
	conn := ConnectionArguments{
		First:  -1,
		Last:   -1,
		Before: "",
		After:  "",
	}
	if filters != nil {
		if first, ok := filters["first"]; ok {
			if first, ok := first.(int); ok {
				conn.First = first
			}
		}
		if last, ok := filters["last"]; ok {
			if last, ok := last.(int); ok {
				conn.Last = last
			}
		}
		if before, ok := filters["before"]; ok {
			conn.Before = ConnectionCursor(fmt.Sprintf("%v", before))
		}
		if after, ok := filters["after"]; ok {
			conn.After = ConnectionCursor(fmt.Sprintf("%v", after))
		}
	}
	return conn
}

func OffsetToCursor(offset int) *string {
	str := fmt.Sprintf("%v%v", PREFIX, offset)
	//return ConnectionCursor(base64.StdEncoding.EncodeToString([]byte(str)))
	return &str
}

// Re-derives the offset from the cursor string.
func CursorToOffset(cursor ConnectionCursor) (int, error) {
	str := string(cursor)
	//str := ""
	//b, err := base64.StdEncoding.DecodeString(string(cursor))
	//if err == nil {
	//	str = string(b)
	//}
	str = strings.Replace(str, PREFIX, "", -1)
	offset, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("Invalid cursor")
	}

	return offset, nil
}

func GetOffsetWithDefault(cursor ConnectionCursor, defaultOffset int) int {
	if cursor == "" {
		return defaultOffset
	}
	offset, err := CursorToOffset(cursor)
	if err != nil {
		return defaultOffset
	}
	return offset
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
