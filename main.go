package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/damilarelana/goGraphQLGoogleAppEngine/middleware"
	"github.com/damilarelana/goGraphQLGoogleAppEngine/resolvers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
)

// Globally initialized Gorilla Mux router
var muxRouter = mux.NewRouter().StrictSlash(true) // instantiate the gorillamux Router and enforce trailing slash rule i.e. `/path` === `/path/`

// Global declaration of schema and err
var schema graphql.Schema // declare GraphQL schema to allow access in other functions
var err error             // declare global error variable

var userType = graphql.NewObject(graphql.ObjectConfig{ // declare GraphQL userType
	Name: "User",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"name":  &graphql.Field{Type: graphql.String},
		"posts": makeListField(makeNodeListType("userTypePostList", postType), resolvers.QueryPostsByUser),
	},
})

var postType = graphql.NewObject(graphql.ObjectConfig{ // declare GraphQL postType
	Name: "Post",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"userID":    &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},
		"content":   &graphql.Field{Type: graphql.String},
	},
})

//
// Mutation
//
var mutationFields = graphql.Fields{ // declare mutation fields: for user, post etc.

	// createUser fields
	"createUser": &graphql.Field{
		Type: userType,
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: resolvers.CreateUser, // call the resolver `createUser`
	},

	// createPost fields
	"createPost": &graphql.Field{
		Type: postType,
		Args: graphql.FieldConfigArgument{
			"userID":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: resolvers.CreatePost, // call the resolver `createPost`
	},
}

var rootMutation = graphql.NewObject(graphql.ObjectConfig{ // declare rootMutation
	Name:   "RootMutation",
	Fields: mutationFields,
})

//
// Query
//
// makeListField function
func makeListField(listType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Field {
	return &graphql.Field{
		Type:    listType,
		Resolve: resolve,
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
	}
}

// makeNodeListType function
func makeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: name,
			Fields: graphql.Fields{
				"nodes":      &graphql.Field{Type: graphql.NewList(nodeType)},
				"totalCount": &graphql.Field{Type: graphql.Int},
			},
		})
}

var rootFields = graphql.Fields{ // declare query fields.
	// queryUser field
	"user": &graphql.Field{
		Type: userType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: resolvers.QueryUser, // call the resolver `queryUser`
	},

	// queryPost field
	"posts": makeListField(makeNodeListType("rootFieldsPostList", postType), resolvers.QueryPosts),
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{ // declare rootQuery
	Name:   "RootQuery",
	Fields: rootFields,
})

// Initialization
// nit builds the schema and maps it to an endpoint handler
func init() {
	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}
	schema, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create a new schema"))
	}
	muxRouter.HandleFunc("/graphql", graphQLHandler)
}

// graphQLServerHomeHandler and entry point for Google App Engine
func graphQLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var stringOutput string

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body) // Read the query via the request body, assuming a POST request
		if err != nil {
			middleware.ResponseError(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		stringOutput = string(body)
	case "GET":
		stringOutput = r.URL.Query().Get("query")
	default:
		stringOutput = ""
	}
	queryParams := graphql.Params{ // compose the GraphQL query parameters
		Schema:        schema,
		RequestString: stringOutput,
		Context:       ctx,
	}

	resp := graphql.Do(queryParams) // execute the GraphQL request

	if len(resp.Errors) > 0 { // check for response errors
		middleware.ResponseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}

	middleware.ResponseJSON(w, resp) // return the query result
}

// Server Home page handler
func graphQLServerHomePageHandler(w http.ResponseWriter, r *http.Request) {
	dataHomePage := "GraphQL Server: homepage"
	io.WriteString(w, dataHomePage)
}

// custom404PageHandler defines custom 404 page
func custom404PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")        // set the content header type
	w.WriteHeader(http.StatusNotFound)                 // this automatically generates a 404 status code
	data404Page := "This page does not exist ... 404!" // page content
	io.WriteString(w, data404Page)
}

func main() {
	muxRouter.NotFoundHandler = http.HandlerFunc(custom404PageHandler) // customer 404 Page handler scenario
	muxRouter.HandleFunc("/", graphQLServerHomePageHandler)
	fmt.Println("GraphQL Server is up and running at http://127.0.0.1:8080")
	for {
		log.Fatal(errors.Wrap(http.ListenAndServe(":8080", muxRouter), "Failed to start GraphQL Server"))
	}
}
