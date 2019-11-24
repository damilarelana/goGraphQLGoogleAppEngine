package resolvers

import (
	"strconv"

	m "github.com/damilarelana/goGraphQLGoogleAppEngine/models"
	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/datastore"
)

// CreateUser function
func CreateUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context

	// Get the name argument
	name, _ := params.Args["name"].(string)
	user := &m.User{Name: name}
	key := datastore.NewIncompleteKey(ctx, "User", nil)

	// Insert user into Datastore
	generatedKey, err := datastore.Put(ctx, key, user)
	if err != nil {
		return m.User{}, err
	}
	user.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	return user, nil
}
