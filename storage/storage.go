package storage

import (
	"context"
	"graphql-ozon/graph/model"
)

type Storage interface {
	CreatePost(ctx context.Context, title string, content string, author string, commentsEnabled bool) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	SetCommentsEnabled(ctx context.Context, postID string, enabled bool, user string) (*model.Post, error)

	CreateComment(ctx context.Context, postID string, parentID *string, text string, author string) (*model.Comment, error)
	GetCommentsByPost(ctx context.Context, postID string) ([]*model.Comment, error)
	GetReplies(ctx context.Context, parentID string) ([]*model.Comment, error)
	GetPostByID(ctx context.Context, postID string) (*model.Post, error)
	GetCommentByID(ctx context.Context, commentID string) (*model.Comment, error)
}
