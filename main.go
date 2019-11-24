package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine"
)

// Entry point for the Google cloud engine
func init() {
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Mutation: rootMutation,
	})
	http.HandleFunc("/", entryPointHandler)
}

// entryPointHandler
func entryPointHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	body, err := ioutil.ReadAll(r.Body) // Read the query
	if err != nil {
		responseError(w, "Invalid request bpdy", http.StatusBadRequest)
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

func main() {
}
