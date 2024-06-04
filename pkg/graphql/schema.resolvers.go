package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database"

	"github.com/Dor1ma/Ozon_Test_Task/internal/database/models"
)

func NewResolver(Repo database.Repository) *Resolver {
	return &Resolver{
		Repo:             Repo,
		CommentObservers: make(map[string]chan *models.Comment),
	}
}
func (r *commentResolver) ID(ctx context.Context, obj *models.Comment) (string, error) {
	return obj.ID, nil
}
func (r *commentResolver) PostID(ctx context.Context, obj *models.Comment) (string, error) {
	return obj.PostID, nil
}
func (r *commentResolver) ParentID(ctx context.Context, obj *models.Comment) (*string, error) {
	return obj.ParentID, nil
}
func (r *postResolver) ID(ctx context.Context, obj *models.Post) (string, error) {
	return obj.ID, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (string, error) {
	return obj.CreatedAt.String(), nil
}

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
	panic("implement me")
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, allowComments bool) (*models.Post, error) {
	return r.Resolver.Repo.CreatePost(title, content, allowComments)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, content string) (*models.Comment, error) {
	comment, err := r.Repo.CreateComment(postID, content, nil)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.CommentObservers[postID]; ok {
		ch <- comment
	}

	return comment, nil
}

// CreateReply is the resolver for the createReply field.
func (r *mutationResolver) CreateReply(ctx context.Context, postID string, content string, parentID string) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: CreateReply - createReply"))
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *models.Post) ([]*models.Comment, error) {
	return r.Resolver.Repo.GetComments(obj.ID)
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*models.Post, error) {
	return r.Resolver.Repo.GetPosts()
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := make(chan *models.Comment, 1)
	r.CommentObservers[postID] = ch

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.CommentObservers, postID)
		r.mu.Unlock()
	}()

	return ch, nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
