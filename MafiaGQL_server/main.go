package main

import (
	"MafiaGQL_server/db"
	"MafiaGQL_server/graph"
	"MafiaGQL_server/graph/generated"
	"MafiaGQL_server/utils"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
)

func main() {
	dbHandle := db.CreateMongoDBHandle()
	err := dbHandle.InitConnection(DB_USERNAME, DB_PASS, DB_HOST, DB_PORT)
	utils.FailOnError("Failed to establish connection to the mongoDB", err)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{dbHandle}}))

	http.Handle("/", playground.Handler("Mafia GraphQL", QUERY_ENDPOINT))
	http.Handle(QUERY_ENDPOINT, srv)

	log.Printf("serving on http://localhost:%s/", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
