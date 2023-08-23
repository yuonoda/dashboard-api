package main

import (
	"dashboard-api/graph"
	"dashboard-api/graph/services"
	"dashboard-api/internal"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"
const dbFile = "./mygraphql.db"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbFile))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("Connected to database %s", dbFile)

	service := services.New(db)

	srv := handler.NewDefaultServer(internal.NewExecutableSchema(internal.Config{Resolvers: &graph.Resolver{
		Srv: service,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
