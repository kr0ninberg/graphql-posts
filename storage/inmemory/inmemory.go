package inmemory

import (
	"context"
	"errors"
	"fmt"
	"graphql-ozon/graph/model"
	"graphql-ozon/storage"
	"strconv"
	"sync"
	"time"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	posts    map[string]*model.Post
	comments map[string]*model.Comment
}

func New() storage.Storage {
	return &InMemoryStorage{
		posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.Comment),
	}
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(len(s.posts) + 1)
	post := &model.Post{
		ID:              id,
		Title:           title,
		Content:         content,
		Author:          author,
		CreatedAt:       time.Now().Format(time.RFC3339),
		CommentsEnabled: commentsEnabled,
	}
	s.posts[id] = post
	return post, nil
}

func (s *InMemoryStorage) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*model.Post, 0, len(s.posts))
	for _, p := range s.posts {
		result = append(result, p)
	}
	return result, nil
}

func (s *InMemoryStorage) SetCommentsEnabled(ctx context.Context, postID string, enabled bool, user string) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	if post.Author != user {
		return nil, errors.New("comment availability change not allowed")
	}
	post.CommentsEnabled = enabled
	return post, nil
}

func (s *InMemoryStorage) CreateComment(ctx context.Context, postID string, parentID *string, text, author string) (*model.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(text) > 2000 {
		return nil, fmt.Errorf("comment text exceeds 2000 characters")
	}

	post, ok := s.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	if !post.CommentsEnabled {
		return nil, errors.New("comments disabled for this post")
	}
	if parentID != nil {
		if _, ok := s.comments[*parentID]; !ok {
			return nil, errors.New("parent comment not found")
		}
	}

	id := strconv.Itoa(len(s.comments) + 1)
	comment := &model.Comment{
		ID:        id,
		PostID:    postID,
		ParentID:  parentID,
		Text:      text,
		Author:    author,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.comments[id] = comment
	return comment, nil
}

func (s *InMemoryStorage) GetCommentsByPost(ctx context.Context, postID string) ([]*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Comment
	for _, c := range s.comments {
		if c.PostID == postID && c.ParentID == nil {
			result = append(result, c)
		}
	}
	return result, nil
}

func (s *InMemoryStorage) GetReplies(ctx context.Context, parentID string) ([]*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var replies []*model.Comment
	for _, c := range s.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			replies = append(replies, c)
		}
	}
	return replies, nil
}

func (s *InMemoryStorage) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *InMemoryStorage) GetCommentByID(ctx context.Context, commentID string) (*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.comments[commentID]
	if !ok {
		return nil, errors.New("comment not found")
	}
	return c, nil
}
