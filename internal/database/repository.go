package database

import "github.com/Dor1ma/Ozon_Test_Task/internal/database/models"

type Repository interface {
	CreatePost(title, content string, allowComments bool) (*models.Post, error)
	GetPosts() ([]*models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	CreateComment(postID int, content string, parentID *int) (*models.Comment, error)
	GetComments(postID int) ([]*models.Comment, error)
}
