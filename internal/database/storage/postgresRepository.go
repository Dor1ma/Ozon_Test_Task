package storage

import (
	"context"
	"github.com/Dor1ma/Ozon_Test_Task/pkg/graphql/model"
	"github.com/jackc/pgx/v4"
	"time"
)

type PostgreSQLRepository struct {
	db *pgx.Conn
}

func NewPostgreSQLRepository(db *pgx.Conn) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) CreatePost(title, content string, allowComments bool) (*model.Post, error) {
	query := `INSERT INTO posts (title, content, allow_comments) VALUES ($1, $2, $3) RETURNING id, title, content, allow_comments`

	var post model.Post
	err := r.db.QueryRow(context.Background(), query, title, content, allowComments).Scan(
		&post.ID, &post.Title, &post.Content, &post.AllowComments,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgreSQLRepository) GetPosts(limit int, after *string) (*model.PostConnection, error) {
	query := `SELECT id, title, content, allow_comments FROM posts ORDER BY id LIMIT $1`
	args := []interface{}{limit}
	if after != nil {
		query = `SELECT id, title, content, allow_comments FROM posts WHERE id > $2 ORDER BY id LIMIT $1`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	edges := make([]*model.PostEdge, len(posts))
	for i, post := range posts {
		edges[i] = &model.PostEdge{
			Cursor: post.ID,
			Node:   post,
		}
	}

	hasNextPage := len(posts) == limit

	return &model.PostConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   &edges[len(edges)-1].Cursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

func (r *PostgreSQLRepository) GetPostByID(id string) (*model.Post, error) {
	query := `SELECT id, title, content, allow_comments FROM posts WHERE id = $1`
	var post model.Post
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.AllowComments,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgreSQLRepository) CreateComment(postID, content string, parentID *string) (*model.Comment, error) {
	query := `INSERT INTO comments (post_id, content, created_at) VALUES ($1, $2, $3) RETURNING id, post_id, content, created_at`

	var comment model.Comment
	err := r.db.QueryRow(context.Background(), query, postID, content, time.Now()).Scan(
		&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if parentID != nil {
		_, err = r.db.Exec(context.Background(), `INSERT INTO replies_comments (parent_comment_id, reply_comment_id) VALUES ($1, $2)`, *parentID, comment.ID)
		if err != nil {
			return nil, err
		}
	}

	return &comment, nil
}

func (r *PostgreSQLRepository) GetComments(postID string, limit int, after *string) (*model.CommentConnection, error) {
	query := `SELECT id, post_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY id LIMIT $2`
	args := []interface{}{postID, limit}
	if after != nil {
		query = `SELECT id, post_id, content, created_at FROM comments WHERE post_id = $1 AND id > $3 ORDER BY id LIMIT $2`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	edges := make([]*model.CommentEdge, len(comments))
	for i, comment := range comments {
		edges[i] = &model.CommentEdge{
			Cursor: comment.ID,
			Node:   comment,
		}
	}

	hasNextPage := len(comments) == limit

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   &edges[len(edges)-1].Cursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

func (r *PostgreSQLRepository) CreateReply(postID, content string, parentID *string) (*model.Comment, error) {
	return r.CreateComment(postID, content, parentID)
}

func (r *PostgreSQLRepository) GetRepliesByCommentID(commentID string, limit int, after *string) (*model.CommentConnection, error) {
	query := `SELECT c.id, c.post_id, c.content, c.created_at FROM comments c JOIN replies_comments rc ON c.id = rc.reply_comment_id WHERE rc.parent_comment_id = $1 ORDER BY c.id LIMIT $2`
	args := []interface{}{commentID, limit}
	if after != nil {
		query = `SELECT c.id, c.post_id, c.content, c.created_at FROM comments c JOIN replies_comments rc ON c.id = rc.reply_comment_id WHERE rc.parent_comment_id = $1 AND c.id > $3 ORDER BY c.id LIMIT $2`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*model.Comment
	for rows.Next() {
		var reply model.Comment
		if err := rows.Scan(&reply.ID, &reply.PostID, &reply.Content, &reply.CreatedAt); err != nil {
			return nil, err
		}
		replies = append(replies, &reply)
	}

	edges := make([]*model.CommentEdge, len(replies))
	for i, reply := range replies {
		edges[i] = &model.CommentEdge{
			Cursor: reply.ID,
			Node:   reply,
		}
	}

	hasNextPage := len(replies) == limit

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   &edges[len(edges)-1].Cursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
