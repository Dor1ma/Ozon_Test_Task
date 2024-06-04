package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database"
	"github.com/Dor1ma/Ozon_Test_Task/pkg/graphql/model"
)

func NewResolver(Repo database.Repository) *Resolver {
	return &Resolver{
		Repo:             Repo,
		CommentObservers: make(map[string]chan *model.Comment),
	}
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, allowComments bool) (*model.Post, error) {
	post, err := r.Repo.CreatePost(title, content, allowComments)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, content string, parentID *string) (*model.Comment, error) {
	comment, err := r.Repo.CreateComment(postID, content, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if ch, ok := r.CommentObservers[postID]; ok {
			ch <- comment
		}
	}()

	return comment, nil
}

// CreateReply is the resolver for the createReply field.
func (r *mutationResolver) CreateReply(ctx context.Context, postID string, parentID string, content string) (*model.Comment, error) {
	comment, err := r.Repo.CreateReply(postID, content, &parentID)
	if err != nil {
		return nil, err
	}

	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if ch, ok := r.CommentObservers[postID]; ok {
			ch <- comment
		}
	}()

	return comment, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, first *int, after *string) (*model.PostConnection, error) {
	limit := 10
	if first != nil {
		limit = *first
	}

	postConnection, err := r.Repo.GetPosts(limit, after)
	if err != nil {
		return nil, err
	}

	return postConnection, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a channel for subscribing to new comments
	ch := make(chan *model.Comment, 1)
	r.CommentObservers[postID] = ch

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.CommentObservers, postID)
	}()

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
