package database

import "github.com/Dor1ma/Ozon_Test_Task/internal/database/models"

type Repository interface {
	CreatePost(title, content string, allowComments bool) (*models.Post, error)
	GetPosts() ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	CreateComment(postID string, content string, parentID *string) (*models.Comment, error)
	GetComments(postID string) ([]*models.Comment, error)
}
