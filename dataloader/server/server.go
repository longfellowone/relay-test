package main

import (
	"log"
	"net/http"

	"dataloader"
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=default password=password dbname=default sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	router := chi.NewRouter()

	router.Handle("/", handler.Playground("Dataloader", "/query"))
	router.Handle("/query", dataloader.LoaderMiddleware(db, handler.GraphQL(
		dataloader.NewExecutableSchema(dataloader.Config{Resolvers: &dataloader.Resolver{DB:db}}),
	)))

	log.Println("connect to http://localhost:8082/ for graphql playground")
	log.Fatal(http.ListenAndServe(":8082", router))
}
