package posts

import "time"

type Posts struct {
	ID        string    `json:"-"`
	UserID    string    `json:"user_ID"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
