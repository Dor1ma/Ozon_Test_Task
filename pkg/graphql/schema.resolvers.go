package graphql

import (
	"context"
	"fmt"
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
	post, err := r.Repo.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %v", err)
	}

	if !post.AllowComments {
		return nil, fmt.Errorf("comments are not allowed for this post")
	}

	if len(content) > 2000 {
		return nil, fmt.Errorf("content too long")
	}

	comment, err := r.Repo.CreateComment(postID, content, parentID)
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
	post, err := r.Repo.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %v", err)
	}

	if !post.AllowComments {
		return nil, fmt.Errorf("comments are not allowed for this post")
	}

	if len(content) > 2000 {
		return nil, fmt.Errorf("content too long")
	}

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

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	post, err := r.Repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, first *int, after *string) (*model.CommentConnection, error) {
	limit := 10
	if first != nil {
		limit = *first
	}

	comments, err := r.Repo.GetComments(postID, limit, after)
	if err != nil {
		return nil, err
	}

	commentEdges := make([]*model.CommentEdge, len(comments.Edges))
	for i, edge := range comments.Edges {
		commentEdges[i] = &model.CommentEdge{
			Cursor: edge.Cursor,
			Node: &model.Comment{
				ID:        edge.Node.ID,
				PostID:    edge.Node.PostID,
				ParentID:  edge.Node.ParentID,
				Content:   edge.Node.Content,
				CreatedAt: edge.Node.CreatedAt,
			},
		}

		// Load replies for the current comment
		replies, err := r.Repo.GetRepliesByCommentID(edge.Node.ID, limit, nil) // nil is temp value to simplify a logic
		if err != nil {
			return nil, err
		}

		// Check if replies exist and have edges
		if replies != nil && len(replies.Edges) > 0 {
			replyEdges := make([]*model.CommentEdge, len(replies.Edges))
			for j, replyEdge := range replies.Edges {
				replyEdges[j] = &model.CommentEdge{
					Cursor: replyEdge.Cursor,
					Node: &model.Comment{
						ID:        replyEdge.Node.ID,
						PostID:    replyEdge.Node.PostID,
						ParentID:  replyEdge.Node.ParentID,
						Content:   replyEdge.Node.Content,
						CreatedAt: replyEdge.Node.CreatedAt,
					},
				}
			}

			commentEdges[i].Node.Replies = &model.CommentConnection{
				Edges:    replyEdges,
				PageInfo: replies.PageInfo,
			}
		} else {
			// If there are no replies, set an empty array
			commentEdges[i].Node.Replies = &model.CommentConnection{
				Edges:    []*model.CommentEdge{},
				PageInfo: &model.PageInfo{HasNextPage: false},
			}
		}
	}

	commentConnection := &model.CommentConnection{
		Edges:    commentEdges,
		PageInfo: comments.PageInfo,
	}

	return commentConnection, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
