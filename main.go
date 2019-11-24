package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
)

// Entry point for the Google cloud engine
func init() {
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Mutation: rootMutation,
	})
}

// entryPointHandler
func graphQLServerHomeHandler(w http.ResponseWriter, r *http.Request) {
	// dataHomePage := "Endpoint: homepage"
	// io.WriteString(w, dataHomePage)
	ctx := appengine.NewContext(r)

	body, err := ioutil.ReadAll(r.Body) // Read the query
	if err != nil {
		responseError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp := graphql.Do(graphql.Params{ // execute the GraphQL request
		Schema:        schema,
		RequestString: string(body),
		Context:       ctx,
	})

	if len(resp.Errors) > 0 { // check for response errors
		responseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}

	responseJSON(w, resp) // return the query result
}

// User type and root mutation
var schema graphql.Schema // declare graphQL schema type

var userType = graphql.NewObject(graphql.ObjectConfig{ // declare graphQL userType
	Name: "User",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{ // declare graphQL
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: createUser,
		},
	},
})

// custom404PageHandler defines custom 404 page
func custom404PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")        // set the content header type
	w.WriteHeader(http.StatusNotFound)                 // this automatically generates a 404 status code
	data404Page := "This page does not exist ... 404!" // page content
	io.WriteString(w, data404Page)
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)                     // instantiate the gorillamux Router and enforce trailing slash rule i.e. `/path` === `/path/`
	muxRouter.NotFoundHandler = http.HandlerFunc(custom404PageHandler) // customer 404 Page handler scenario
	muxRouter.HandleFunc("/", graphQLServerHomeHandler)
	fmt.Println("GraphQL Server is up and running at http://127.0.0.1:8080")
	for {
		log.Fatal(errors.Wrap(http.ListenAndServe(":8080", muxRouter), "Failed to start GraphQL Server"))
	}
}
