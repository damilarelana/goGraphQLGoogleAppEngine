package resolvers

import (
	"context"
	"errors"
	"strconv"
	"time"

	m "github.com/damilarelana/goGraphQLGoogleAppEngine/models"
	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/datastore"
)

// PostListResult struct
type PostListResult struct {
	Nodes      []m.Post `json:"nodes"`
	TotalCount int      `json:"totalCount"`
}

func queryPostList(ctx context.Context, query *datastore.Query) (PostListResult, error) {
	query = query.Order("-CreatedAt") // order by creation time
	var result PostListResult
	keys, err := query.GetAll(ctx, &result.Nodes) // run the query
	if err != nil {
		return result, err
	} else {
		for i, key := range keys { // set IDs
			result.Nodes[i].ID = strconv.FormatInt(key.IntID(), 10)
		}
		result.TotalCount = len(result.Nodes) // set total count
	}
	return result, nil
}

// QueryPosts function
func QueryPosts(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("Post")
	limit, ok := params.Args["limit"].(int)
	if ok {
		query = query.Limit(limit)
	}
	offset, ok := params.Args["offset"].(int)
	if ok {
		query = query.Offset(offset)
	}
	return queryPostList(ctx, query)
}

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

// CreatePost function
func CreatePost(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context

	// Get the arguments
	content, _ := params.Args["content"].(string)
	userID, _ := params.Args["userID"].(string)
	post := &m.Post{UserID: userID, Content: content, CreatedAt: time.Now().UTC()}
	key := datastore.NewIncompleteKey(ctx, "Post", nil)

	// Insert post into Datastore
	generatedKey, err := datastore.Put(ctx, key, post)
	if err != nil {
		return m.Post{}, err
	}
	post.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	return post, nil
}

// QueryUser function
func QueryUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context

	strID, ok := params.Args["id"].(string)
	if ok {
		id, err := strconv.ParseInt(strID, 10, 64) // Parse ID argument
		if err != nil {
			return nil, errors.New("Invalid id")
		}
		user := &m.User{ID: strID}
		key := datastore.NewKey(ctx, "User", "", id, nil)

		err = datastore.Get(ctx, key, user) // Fetch user by ID
		if err != nil {
			return nil, errors.New("User not found")
		}
		return user, nil
	}
	return m.User{}, nil
}

// QueryPostsByUser function
func QueryPostsByUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("Post")

	limit, ok := params.Args["limit"].(int)
	if ok {
		query = query.Limit(limit)
	}
	offset, ok := params.Args["offset"].(int)
	if ok {
		query = query.Offset(offset)
	}

	// check user's ID against post's UserID field
	user, ok := params.Source.(*m.User)
	if ok {
		query = query.Filter("UserID =", user.ID)
	}
	return queryPostList(ctx, query)
}
