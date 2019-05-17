////go:generate go run github.com/vektah/dataloaden OrderConnectionLoader string *dataloader.OrderConnection
////go:generate go run github.com/vektah/dataloaden ProjectConnectionLoader string *dataloader.ProjectConnection
//go:generate go run github.com/vektah/dataloaden OrderIDsByProjectLoader string []string

package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

type loaders struct {
	orderIDsByProject *OrderIDsByProjectLoader
}

func LoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}

		wait := 250 * time.Microsecond

		ldrs.orderIDsByProject = &OrderIDsByProjectLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []string) ([][]string, []error) {
				fmt.Printf("SELECT * FROM orders WHERE project_id IN (%s)\n", strings.Join(keys, ","))
				time.Sleep(5 * time.Millisecond)

				orderIDs := make([][]string, len(keys))
				errors := make([]error, len(keys))
				for i := range keys {
					orderIDs[i] = []string{"1", "2", "3"}
				}

				return orderIDs, errors
			},
		}

		//ldrs.ordersByProject = &OrderConnectionLoader{
		//	wait:     wait,
		//	maxBatch: 100,
		//	fetch: func(keys []string) ([]*OrderConnection, []error) {
		//		fmt.Printf("SELECT * FROM orders WHERE project_id IN (%s)\n", strings.Join(keys, ","))
		//		time.Sleep(5 * time.Millisecond)
		//
		//		orderConnections := make([]*OrderConnection, len(keys))
		//		errors := make([]error, len(keys))
		//		for i, key := range keys {
		//
		//			edges := []*OrderEdge{
		//				{Node: &Order{ID: key}, Cursor: "testing"},
		//			}
		//
		//			orderConnections[i] = &OrderConnection{Edges: edges, PageInfo: &PageInfo{HasPreviousPage: true}}
		//		}
		//		return orderConnections, errors
		//	},
		//}

		dlCtx := context.WithValue(r.Context(), ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(dlCtx))
	})
}

func ctxLoaders(ctx context.Context) loaders {
	return ctx.Value(ctxKey).(loaders)
}
