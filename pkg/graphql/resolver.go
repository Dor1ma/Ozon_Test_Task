package graphql

import (
	"github.com/Dor1ma/Ozon_Test_Task/internal/database"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/models"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo             database.Repository
	CommentObservers map[string]chan *models.Comment
	mu               sync.Mutex
}
