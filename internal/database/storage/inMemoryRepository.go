package storage

import (
	"errors"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/models"
	"sync"
	"time"
)

// Имплементация репозитория для работы в режиме "in_memory"

type InMemoryRepository struct {
	posts     map[int]*models.Post
	comments  map[int][]*models.Comment
	postIDSeq int // Счетчик для генерации уникальных ID постов
	mutex     sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts:    make(map[int]*models.Post),
		comments: make(map[int][]*models.Comment),
	}
}

func (r *InMemoryRepository) CreatePost(title, content string, allowComments bool) (*models.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	post := &models.Post{
		ID:            r.postIDSeq,
		Title:         title,
		Content:       content,
		AllowComments: allowComments,
	}

	r.posts[post.ID] = post
	r.postIDSeq++

	return post, nil
}

func (r *InMemoryRepository) GetPosts() ([]*models.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	posts := make([]*models.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *InMemoryRepository) GetPostById(id int) (*models.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (r *InMemoryRepository) CreateComment(postID int, content string, parentID *int) (*models.Comment, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}

	comment := &models.Comment{
		ID:        len(r.comments[postID]) + 1,
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	r.comments[postID] = append(r.comments[postID], comment)

	return comment, nil
}

func (r *InMemoryRepository) GetComments(postID int) ([]*models.Comment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	comments, ok := r.comments[postID]
	if !ok {
		return nil, errors.New("comments not found")
	}

	return comments, nil
}
