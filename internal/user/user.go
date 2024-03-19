package user

import "time"

type User struct {
	ID             string
	Name           string
	Email          *string
	PhoneNumber    *string
	FriendCount    int
	ImageURL       *string
	HashedPassword string
	CreatedAt      time.Time
}
