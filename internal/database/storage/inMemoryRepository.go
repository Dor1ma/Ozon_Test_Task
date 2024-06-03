package storage

import (
	"errors"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/models"
	"strconv"
	"sync"
	"time"
)

// Имплементация репозитория для работы в режиме "in_memory"

type InMemoryRepository struct {
	posts     map[string]*models.Post
	comments  map[string][]*models.Comment
	postIDSeq int // Счетчик для генерации уникальных ID постов
	mutex     sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts:    make(map[string]*models.Post),
		comments: make(map[string][]*models.Comment),
	}
}

func (r *InMemoryRepository) CreatePost(title, content string, allowComments bool) (*models.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	post := &models.Post{
		ID:            strconv.Itoa(r.postIDSeq),
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

func (r *InMemoryRepository) GetPostByID(id string) (*models.Post, error) { // изменено: параметр и возвращаемое значение теперь string
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (r *InMemoryRepository) CreateComment(postID string, content string, parentID *string) (*models.Comment, error) { // изменено: параметр теперь string
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}

	comment := &models.Comment{
		ID:        strconv.Itoa(len(r.comments[postID]) + 1),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	r.comments[postID] = append(r.comments[postID], comment)

	return comment, nil
}

func (r *InMemoryRepository) GetComments(postID string) ([]*models.Comment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	comments, ok := r.comments[postID]
	if !ok {
		return nil, errors.New("comments not found")
	}

	return comments, nil
}
