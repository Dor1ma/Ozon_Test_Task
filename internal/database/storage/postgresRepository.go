package storage

import (
	"context"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/models"
	"github.com/jackc/pgx/v4"
	"time"
)

// Имплементация репозитория для работы в режиме "postgres"

type PostgreSQLRepository struct {
	db *pgx.Conn
}

func NewPostgreSQLRepository(db *pgx.Conn) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) CreatePost(title string, content string, allowComments bool) (*models.Post, error) {
	query :=
		`INSERT INTO posts (title, content, allow_comments)
	VALUES ($1, $2, $3)
	RETURNING id, title, content, allow_comments`

	var post models.Post
	resultChan := make(chan error)

	go func() {
		err := r.db.QueryRow(context.Background(), query, title, content, allowComments).Scan(
			&post.ID, &post.Title, &post.Content, &post.AllowComments,
		)
		resultChan <- err
	}()

	err := <-resultChan
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgreSQLRepository) GetPosts() ([]*models.Post, error) {
	query := `
        SELECT id, title, content, allow_comments
        FROM posts
    `
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostgreSQLRepository) GetPostByID(id string) (*models.Post, error) {
	query := `
        SELECT id, title, content, allow_comments
        FROM posts
        WHERE id = $1
    `
	var post models.Post
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.AllowComments,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgreSQLRepository) CreateComment(postID string, content string, parentID *string) (*models.Comment, error) {
	query := `
		INSERT INTO comments (post_id, content, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, content, created_at
	`

	var comment models.Comment
	err := r.db.QueryRow(context.Background(), query, postID, content, parentID, time.Now()).Scan(
		&comment.ID, &comment.PostID, &comment.Content, &comment.ParentID, &comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *PostgreSQLRepository) GetComments(postID string) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostgreSQLRepository) CreateReply(postID string, content string, parentID *string) (*models.Comment, error) {
	_, err := r.CreateComment(postID, content, parentID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO replies_comments (parent_comment_id, reply_comment_id)
		VALUES ($1, $2)
		RETURNING parent_comment_id
		`

	var comment models.Comment
	err = r.db.QueryRow(context.Background(), query, postID, content, parentID, time.Now()).Scan(
		&comment.ID, &comment.PostID, &comment.Content, &comment.ParentID, &comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *PostgreSQLRepository) GetRepliesByPostID(postID string) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, rc.parent_comment_id, c.content, c.created_at
		FROM comments c
		JOIN replies_comments rc ON c.id = rc.reply_comment_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC	
		`

	rows, err := r.db.Query(context.Background(), query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
