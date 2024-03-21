package posts

import "time"

type Posts struct {
	ID        string
	UserID    string
	Content   string
	Tags      []string
	CreatedAt time.Time
}

type Comment struct {
	ID        uint64
	UserID    string
	PostID    string
	Content   string
	CreatedAt time.Time
}
