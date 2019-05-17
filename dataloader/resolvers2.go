//go:generate go run github.com/99designs/gqlgen

package dataloader

import (
	"context"
	"errors"
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
	IDs   []string
	After *string
	First *int
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
	//fmt.Println("he")
	//for _, value := range obj.IDs {
	//	fmt.Println(value)
	//}

	//friends, err := r.resolveCharacters(ctx, obj.Ids)
	//if err != nil {
	//	return nil, err
	//}
	//
	//edges := make([]*models.FriendsEdge, obj.To-obj.From)
	//for i := range edges {
	//	edges[i] = &models.FriendsEdge{
	//		Cursor: models.EncodeCursor(obj.From + i),
	//		Node:   friends[obj.From+i],
	//	}
	//}
	//return edges, nil

	//fmt.Println(obj.IDs)
	return []*OrderEdge{}, nil
}

func (r *Resolver) resolveOrders(ctx context.Context, ids []string) ([]*Order, error) {
	//var result []models.Character
	//for _, id := range ids {
	//	char, err := r.Query().Character(ctx, id)
	//	if err != nil {
	//		return nil, err
	//	}
	//	result = append(result, char)
	//}
	//return result, nil
	return []*Order{}, nil

}

type ArraySliceMetaInfo struct {
	SliceStart  int `json:"sliceStart"`
	ArrayLength int `json:"arrayLength"`
}

type ConnectionCursor string

type ConnectionArguments struct {
	Before ConnectionCursor `json:"before"`
	After  ConnectionCursor `json:"after"`
	First  int              `json:"first"` // -1 for undefined, 0 would return zero results
	Last   int              `json:"last"`  //  -1 for undefined, 0 would return zero results
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

func (r *Resolver) resolveOrderConnection(orderIDs []string, after *string, first *int, before *string, last *int) (*OrderConnection, error) {

	fmt.Println("-------------------------------------------------")
	filter := map[string]interface{}{
		"first": 2,
	}
	args := NewConnectionArguments(filter)

	meta := ArraySliceMetaInfo{
		SliceStart:  0,
		ArrayLength: len(orderIDs),
	}

	arraySlice := orderIDs

	sliceEnd := meta.SliceStart + len(arraySlice)
	beforeOffset := GetOffsetWithDefault(args.Before, meta.ArrayLength)
	afterOffset := GetOffsetWithDefault(args.After, -1)

	startOffset := ternaryMax(meta.SliceStart-1, afterOffset, -1) + 1
	endOffset := ternaryMin(sliceEnd, beforeOffset, meta.ArrayLength)

	if args.First != -1 {
		endOffset = min(endOffset, startOffset+args.First)
	}

	if args.Last != -1 {
		startOffset = max(startOffset, endOffset-args.Last)
	}

	begin := max(startOffset-meta.SliceStart, 0)
	end := len(arraySlice) - (sliceEnd - endOffset)

	if begin > end {
		return &OrderConnection{}, nil
	}

	slice := arraySlice[begin:end]

	for index, value := range slice {
		fmt.Println(value, OffsetToCursor(startOffset+index))
	}

	var firstEdgeCursor, lastEdgeCursor ConnectionCursor
	if len(slice) > 0 {
		firstEdgeCursor = OffsetToCursor(startOffset)
		lastEdgeCursor = OffsetToCursor(startOffset + (len(slice) - 1))
	}

	lowerBound := 0
	if len(args.After) > 0 {
		lowerBound = afterOffset + 1
	}

	upperBound := meta.ArrayLength
	if len(args.Before) > 0 {
		upperBound = beforeOffset
	}

	hasPreviousPage := false
	if args.Last != -1 {
		hasPreviousPage = startOffset > lowerBound
	}

	hasNextPage := false
	if args.First != -1 {
		hasNextPage = endOffset < upperBound
	}

	fmt.Println(firstEdgeCursor, lastEdgeCursor, hasNextPage, hasPreviousPage)

	return &OrderConnection{
		IDs:   orderIDs,
		After: after,
		First: first,
	}, nil
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

	return []*Project{
		{ID: "projectID01"},
		{ID: "projectID02"},
	}, nil
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
const PREFIX = "cursor:"

func OffsetToCursor(offset int) ConnectionCursor {
	str := fmt.Sprintf("%v%v", PREFIX, offset)
	//return ConnectionCursor(base64.StdEncoding.EncodeToString([]byte(str)))
	return ConnectionCursor(str)
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

func ternaryMax(a, b, c int) int {
	return max(max(a, b), c)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func ternaryMin(a, b, c int) int {
	return min(min(a, b), c)
}
