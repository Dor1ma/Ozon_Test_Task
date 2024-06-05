package tests

import (
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/mocks"
	"testing"

	"github.com/Dor1ma/Ozon_Test_Task/internal/database/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)

	repo := storage.NewPostgreSQLRepository(mockDB)

	authorID := "author123"
	title := "Test Post"
	content := "Lorem ipsum"
	allowComments := true

	mockRow := mocks.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(
		gomock.Any(), // id
		gomock.Any(), // author_id
		gomock.Any(), // title
		gomock.Any(), // content
		gomock.Any(), // allow_comments
	).Return(nil).Times(1)

	mockDB.EXPECT().
		QueryRow(gomock.Any(), gomock.Any(), authorID, title, content, allowComments).
		Return(mockRow).
		Times(1)

	_, err := repo.CreatePost(authorID, title, content, allowComments)

	assert.NoError(t, err, "Expected no error")
}
