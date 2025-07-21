package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"graphql-ozon/graph/model"
	"graphql-ozon/storage"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func New(dsn string) (storage.Storage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*model.Post, error) {
	var post model.Post
	err := s.db.QueryRowxContext(ctx,
		`INSERT INTO posts (title, content, author, comments_enabled)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, title, content, author, created_at, comments_enabled`,
		title, content, author, commentsEnabled).StructScan(&post)
	if err != nil {
		return nil, err
	}
	post.ID = fmt.Sprintf("%d", post.ID) // convert from int to string
	return &post, nil
}

func (s *PostgresStorage) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post
	err := s.db.SelectContext(ctx, &posts,
		`SELECT id, title, content, author, created_at, comments_enabled FROM posts ORDER BY id`)
	for _, p := range posts {
		p.ID = fmt.Sprintf("%v", p.ID)
	}
	return posts, err
}

func (s *PostgresStorage) SetCommentsEnabled(ctx context.Context, postID string, enabled bool, user string) (*model.Post, error) {
	var post model.Post
	err := s.db.GetContext(ctx, &post,
		`SELECT id, title, content, author, created_at, comments_enabled FROM posts WHERE id = $1`, postID)
	if err != nil {
		return nil, err
	}

	if post.Author != user {
		return nil, errors.New("you are not the author of this post")
	}

	err = s.db.QueryRowxContext(ctx,
		`UPDATE posts SET comments_enabled = $2 WHERE id = $1
		 RETURNING id, title, content, author, created_at, comments_enabled`,
		postID, enabled).StructScan(&post)
	if err != nil {
		return nil, err
	}
	post.ID = fmt.Sprintf("%v", post.ID)
	return &post, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, postID string, parentID *string, text, author string) (*model.Comment, error) {
	if len(text) > 2000 {
		return nil, errors.New("text exceeds 2000 characters")
	}

	var comment model.Comment
	err := s.db.QueryRowxContext(ctx,
		`INSERT INTO comments (post_id, parent_id, text, author)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, post_id, parent_id, text, author, created_at`,
		postID, parentID, text, author).StructScan(&comment)
	if err != nil {
		return nil, err
	}
	comment.ID = fmt.Sprintf("%v", comment.ID)
	comment.PostID = fmt.Sprintf("%v", comment.PostID)
	if comment.ParentID != nil {
		id := fmt.Sprintf("%v", *comment.ParentID)
		comment.ParentID = &id
	}
	return &comment, nil
}

func (s *PostgresStorage) GetCommentsByPost(ctx context.Context, postID string) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := s.db.SelectContext(ctx, &comments,
		`SELECT id, post_id, parent_id, text, author, created_at
		 FROM comments
		 WHERE post_id = $1 AND parent_id IS NULL
		 ORDER BY id`, postID)
	for _, c := range comments {
		c.ID = fmt.Sprintf("%v", c.ID)
		c.PostID = fmt.Sprintf("%v", c.PostID)
		if c.ParentID != nil {
			id := fmt.Sprintf("%v", *c.ParentID)
			c.ParentID = &id
		}
	}
	return comments, err
}

func (s *PostgresStorage) GetReplies(ctx context.Context, parentID string) ([]*model.Comment, error) {
	var replies []*model.Comment
	err := s.db.SelectContext(ctx, &replies,
		`SELECT id, post_id, parent_id, text, author, created_at
		 FROM comments
		 WHERE parent_id = $1
		 ORDER BY id`, parentID)
	for _, c := range replies {
		c.ID = fmt.Sprintf("%v", c.ID)
		c.PostID = fmt.Sprintf("%v", c.PostID)
		if c.ParentID != nil {
			id := fmt.Sprintf("%v", *c.ParentID)
			c.ParentID = &id
		}
	}
	return replies, err
}

func (s *PostgresStorage) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	var post model.Post
	err := s.db.GetContext(ctx, &post,
		`SELECT id, title, content, author, created_at, comments_enabled FROM posts WHERE id = $1`, postID)
	if err != nil {
		return nil, err
	}
	post.ID = fmt.Sprintf("%v", post.ID)
	return &post, nil
}

func (s *PostgresStorage) GetCommentByID(ctx context.Context, commentID string) (*model.Comment, error) {
	var c model.Comment
	err := s.db.GetContext(ctx, &c,
		`SELECT id, post_id, parent_id, text, author, created_at FROM comments WHERE id = $1`, commentID)
	if err != nil {
		return nil, err
	}
	c.ID = fmt.Sprintf("%v", c.ID)
	c.PostID = fmt.Sprintf("%v", c.PostID)
	if c.ParentID != nil {
		id := fmt.Sprintf("%v", *c.ParentID)
		c.ParentID = &id
	}
	return &c, nil
}
