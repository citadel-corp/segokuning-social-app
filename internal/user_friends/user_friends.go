package userfriends

import "time"

type UserFriends struct {
	ID        uint64
	UserID    string
	FriendID  string
	CreatedAt time.Time
}
