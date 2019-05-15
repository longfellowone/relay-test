package server

//
//import (
//	"context"
//	"fmt"
//) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
//
//type Resolver struct{}
//
//func (r *Resolver) Query() QueryResolver {
//	return &queryResolver{r}
//}
//func (r *Resolver) Order() OrderResolver {
//	return &orderResolver{r}
//}
//
//
//type queryResolver struct{ *Resolver }
//
//func (r *queryResolver) Orders(ctx context.Context, projectId *string, after *string, first *int, before *string, last *int) (*OrderConnection, error) {
//	fmt.Println("queryResolver/Orders")
//	fmt.Println(projectId)
//
//	order1 := &OrderEdge{
//		Node: &Order{
//			ID:        "Testing1",
//			Timestamp: "",
//		},
//		Cursor: "",
//	}
//
//	order2 := &OrderEdge{
//		Node: &Order{
//			ID:        "Testing2",
//			Timestamp: "",
//		},
//		Cursor: "",
//	}
//
//	edges := []*OrderEdge{order1, order2}
//
//	return &OrderConnection{
//		PageInfo: PageInfo{
//			HasNextPage:     true,
//			HasPreviousPage: false,
//			StartCursor:     nil,
//			EndCursor:       nil,
//		},
//		Edges: edges,
//	}, nil
//}
//func (r *queryResolver) Node(ctx context.Context, id string) (Node, error) {
//	panic("not implemented")
//}
//
//type orderResolver struct{ *Resolver }
//
//func (r *orderResolver) Project(ctx context.Context, obj *Order) (*Project, error) {
//	fmt.Println("orderResolver/Project")
//
//	return &Project{
//		ID:   "testing id",
//		Name: "Project name",
//	}, nil
//}
//
////type orderConnectionResolver struct{ *Resolver }
//
////func (r *orderConnectionResolver) PageInfo(ctx context.Context, obj *OrderConnection) (*PageInfo, error) {
////	fmt.Println("orderConnectionResolver/PageInfo")
////
////	return &PageInfo{}, nil
////}
////func (r *orderConnectionResolver) Edges(ctx context.Context, obj *OrderConnection) ([]*OrderEdge, error) {
////	fmt.Println("orderConnectionResolver/Edges")
////
////	return nil, nil
////}
