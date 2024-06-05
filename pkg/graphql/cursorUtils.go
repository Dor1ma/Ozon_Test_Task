package graphql

import (
	"encoding/base64"
	"fmt"
	"github.com/Dor1ma/Ozon_Test_Task/pkg/graphql/model"
)

func encodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Post:%s", id)))
}

func getEndCursor(posts []*model.Post) *string {
	if len(posts) == 0 {
		return nil
	}
	cursor := encodeCursor(posts[len(posts)-1].ID)
	return &cursor
}
