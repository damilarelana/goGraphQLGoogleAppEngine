package main

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/datastore"
)

func createUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context

	// Get the name argument
	name, _ := params.Args["name"].(string)
	user := &User{Name: name}
	key := datastore.NewIncompleteKey(ctx, "User", nil)

	// Insert user into Datastore
	generatedKey, err := datastore.Put(ctx, key, user)
	if err != nil {
		return User{}, err
	}
	user.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	return user, nil
}
