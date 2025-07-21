// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package graph

import "graphql-ozon/graph/model"

type Resolver struct {
	PostsContainer    map[string]*model.Post
	CommentsContainer map[string]*model.Comment
}
