////go:generate go run github.com/vektah/dataloaden OrderConnectionLoader string *dataloader.OrderConnection
////go:generate go run github.com/vektah/dataloaden ProjectConnectionLoader string *dataloader.ProjectConnection
//go:generate go run github.com/vektah/dataloaden OrderIDsByProjectLoader string []string
//go:generate go run github.com/vektah/dataloaden OrderLoader string *dataloader.Order

package dataloader

import (
	"context"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

type loaders struct {
	orderIDsByProject *OrderIDsByProjectLoader
	orderLoader       *OrderLoader
}

func LoaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}

		wait := 250 * time.Microsecond

		ldrs.orderIDsByProject = &OrderIDsByProjectLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []string) ([][]string, []error) {
				time.Sleep(5 * time.Millisecond)
				errors := make([]error, len(keys))

				query, args, err := sqlx.In("SELECT orderid,projectid FROM orders WHERE orders.projectid IN (?) ORDER by sentdate DESC", keys)
				query = db.Rebind(query)

				rows, err := db.Query(query, args...)
				if err != nil {
					errors = append(errors, err)
				}

				orderIDs := make(map[string][]string)

				for rows.Next() {
					var orderid string
					var projectid string

					if err := rows.Scan(&orderid, &projectid); err != nil {
						errors = append(errors, err)
					}

					orderIDs[projectid] = append(orderIDs[projectid], orderid)
				}

				projectOrderIDs := make([][]string, len(keys))
				for i, v := range keys {
					projectOrderIDs[i] = orderIDs[v]
				}

				return projectOrderIDs, errors
			},
		}

		ldrs.orderLoader = &OrderLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []string) ([]*Order, []error) {
				time.Sleep(5 * time.Millisecond)

				errors := make([]error, len(keys))

				query, args, err := sqlx.In("SELECT orderid,projectid FROM orders WHERE orderid IN (?) ORDER by sentdate DESC", keys)
				query = db.Rebind(query)

				rows, err := db.Query(query, args...)
				if err != nil {
					errors = append(errors, err)
				}

				orders := make([]*Order, 0)

				var orderid string
				var projectid string

				for rows.Next() {
					if err := rows.Scan(&orderid, &projectid); err != nil {
						errors = append(errors, err)
					}
					orders = append(orders, &Order{ID: orderid, ProjectId: projectid})
				}

				return orders, errors
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
