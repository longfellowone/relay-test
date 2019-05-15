//go:generate go run github.com/vektah/dataloaden -keys int dataloader.Address
//go:generate go run github.com/vektah/dataloaden -keys int -slice dataloader.Order
//go:generate go run github.com/vektah/dataloaden -keys int dataloader.OrderConnection
//go:generate go run github.com/vektah/dataloaden -keys int -slice dataloader.Item

// https://medium.freecodecamp.org/deep-dive-into-graphql-with-golang-d3e02a429ac3

package dataloader

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

type loaders struct {
	ordersByCustomer *OrderSliceLoader
	ordersByProject  *OrderConnectionLoader
	//addressByID      *AddressLoader
	//itemsByOrder     *ItemSliceLoader
}

func LoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}

		wait := 250 * time.Microsecond

		// 1:M loader
		ldrs.ordersByCustomer = &OrderSliceLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([][]Order, []error) {
				var keySql []string
				for _, key := range keys {
					keySql = append(keySql, strconv.Itoa(key))
				}

				fmt.Printf("SELECT * FROM orders WHERE customer_id IN (%s)\n", strings.Join(keySql, ","))
				time.Sleep(5 * time.Millisecond)

				orders := make([][]Order, len(keys))
				errors := make([]error, len(keys))
				for i, key := range keys {
					id := 10 + rand.Int()%3
					orders[i] = []Order{
						{ID: id, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
						{ID: id + 1, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
					}

					// if you had another customer loader you would prime its cache here
					// by calling `ldrs.ordersByID.Prime(id, orders[i])`
				}

				return orders, errors
			},
		}

		// 1:1 loader
		ldrs.ordersByProject = &OrderConnectionLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([]*OrderConnection, []error) {
				var keySql []string
				for _, key := range keys {
					keySql = append(keySql, strconv.Itoa(key))
				}

				fmt.Printf("SELECT * FROM orders WHERE project_id IN (%s)\n", strings.Join(keySql, ","))
				time.Sleep(5 * time.Millisecond)

				// make edge array
				orderConnections := make([]*OrderConnection, len(keys))
				errors := make([]error, len(keys))
				for i, key := range keys {

					// range over [i]
					edges := []*OrderEdge{
						{Node: &Order{ID: key, Date: time.Now()}, Cursor: "testing"},
					}

					orderConnections[i] = &OrderConnection{Edges: edges, PageInfo: PageInfo{HasPreviousPage: true}}
				}
				return orderConnections, errors
			},
		}

		// simple 1:1 loader, fetch an address by its primary key
		//ldrs.addressByID = &AddressLoader{
		//	wait:     wait,
		//	maxBatch: 100,
		//	fetch: func(keys []int) ([]*Address, []error) {
		//		var keySql []string
		//		for _, key := range keys {
		//			keySql = append(keySql, strconv.Itoa(key))
		//		}
		//
		//		fmt.Printf("SELECT * FROM address WHERE id IN (%s)\n", strings.Join(keySql, ","))
		//		time.Sleep(5 * time.Millisecond)
		//
		//		addresses := make([]*Address, len(keys))
		//		errors := make([]error, len(keys))
		//		for i, key := range keys {
		//			addresses[i] = &Address{ID: key, Street: "home street", Country: "hometon " + strconv.Itoa(key)}
		//		}
		//		return addresses, errors
		//	},
		//}

		//// 1:M loader
		//ldrs.ordersByProject = &OrderConnectionSliceLoader{
		//	wait:     wait,
		//	maxBatch: 100,
		//	fetch: func(keys []int) ([][]OrderConnection, []error) {
		//		var keySql []string
		//		for _, key := range keys {
		//			keySql = append(keySql, strconv.Itoa(key))
		//		}
		//
		//		fmt.Printf("SELECT * FROM orders WHERE project_id IN (%s)\n", strings.Join(keySql, ","))
		//		time.Sleep(5 * time.Millisecond)
		//
		//		orders := make([][]OrderConnection, len(keys))
		//		errors := make([]error, len(keys))
		//		for i, key := range keys {
		//			orders[i] = []OrderConnection{
		//				{Edges: []*OrderEdge{{Node: &Order{ID: key}}}, PageInfo: PageInfo{}},
		//				//{ID: id, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
		//				//{ID: id + 1, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
		//			}
		//
		//			// if you had another customer loader you would prime its cache here
		//			// by calling `ldrs.ordersByID.Prime(id, orders[i])`
		//		}
		//
		//		return orders, errors
		//	},
		//}

		//// M:M loader
		//ldrs.itemsByOrder = &ItemSliceLoader{
		//	wait:     wait,
		//	maxBatch: 100,
		//	fetch: func(keys []int) ([][]Item, []error) {
		//		var keySql []string
		//		for _, key := range keys {
		//			keySql = append(keySql, strconv.Itoa(key))
		//		}
		//
		//		fmt.Printf("SELECT * FROM items JOIN item_order WHERE item_order.order_id IN (%s)\n", strings.Join(keySql, ","))
		//		time.Sleep(5 * time.Millisecond)
		//
		//		items := make([][]Item, len(keys))
		//		errors := make([]error, len(keys))
		//		for i := range keys {
		//			items[i] = []Item{
		//				{Name: "item " + strconv.Itoa(rand.Int()%20+20)},
		//				{Name: "item " + strconv.Itoa(rand.Int()%20+20)},
		//			}
		//		}
		//
		//		return items, errors
		//	},
		//}

		dlCtx := context.WithValue(r.Context(), ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(dlCtx))
	})
}

func ctxLoaders(ctx context.Context) loaders {
	return ctx.Value(ctxKey).(loaders)
}
