package inmemory_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"graphql-ozon/storage/inmemory"
	"testing"
)

func TestCreateAndGetPosts(t *testing.T) {
	store := inmemory.New()
	ctx := context.Background()

	post, err := store.CreatePost(ctx, "Test Title", "Test Content", "test author", true)
	require.NoError(t, err)
	require.Equal(t, "Test Title", post.Title)

	posts, err := store.GetAllPosts(ctx)
	require.NoError(t, err)
	require.Len(t, posts, 1)
	assert.Equal(t, post.ID, posts[0].ID)
}

func TestCreateCommentAndGetReplies(t *testing.T) {
	store := inmemory.New()
	ctx := context.Background()

	post, _ := store.CreatePost(ctx, "Post", "Content", "test author", true)

	comment, err := store.CreateComment(ctx, post.ID, nil, "root", "root-comment")
	require.NoError(t, err)
	require.Equal(t, "root", comment.Text)

	child, err := store.CreateComment(ctx, post.ID, &comment.ID, "child", "child-comment")
	require.NoError(t, err)
	require.Equal(t, comment.ID, *child.ParentID)

	replies, err := store.GetReplies(ctx, comment.ID)
	require.NoError(t, err)
	require.Len(t, replies, 1)
	assert.Equal(t, "child", replies[0].Text)
}

func TestSetCommentsEnabled(t *testing.T) {
	store := inmemory.New()
	ctx := context.Background()

	post, _ := store.CreatePost(ctx, "Post", "Content", "test author", true)

	_, err := store.SetCommentsEnabled(ctx, post.ID, false, "someone_else")
	require.Error(t, err)

	updated, err := store.SetCommentsEnabled(ctx, post.ID, false, post.Author)
	require.NoError(t, err)
	assert.False(t, updated.CommentsEnabled)
}

func TestCreateComment_Disabled(t *testing.T) {
	store := inmemory.New()
	ctx := context.Background()

	post, _ := store.CreatePost(ctx, "Post", "Content", "test author", false)

	_, err := store.CreateComment(ctx, post.ID, nil, "user", "no comments allowed")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "comments disabled")
}
