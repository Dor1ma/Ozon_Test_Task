package database

import "github.com/Dor1ma/Ozon_Test_Task/pkg/graphql/model"

type Repository interface {
	CreatePost(title, content string, allowComments bool) (*model.Post, error)
	GetPosts(limit int, after *string) (*model.PostConnection, error)
	GetPostByID(id string) (*model.Post, error)
	CreateComment(postID string, content string, parentID *string) (*model.Comment, error)
	GetComments(postID string, limit int, after *string) (*model.CommentConnection, error)
	CreateReply(postID string, content string, parentID *string) (*model.Comment, error)
	GetRepliesByCommentID(commentID string, limit int, after *string) (*model.CommentConnection, error)
}
