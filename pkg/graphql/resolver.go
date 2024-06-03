package graphql

import "github.com/Dor1ma/Ozon_Test_Task/internal/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repo database.Repository
}
