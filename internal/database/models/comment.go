package models

import "time"

type Comment struct {
	ID        int
	PostID    int
	ParentID  *int
	Content   string
	CreatedAt time.Time
}
