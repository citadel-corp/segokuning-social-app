package user

import "time"

type UserRegisterResponse struct {
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Name        string  `json:"name"`
	AccessToken string  `json:"accessToken"`
}

type UserLoginResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserListResponse struct {
	ID          string    `json:"userId"`
	Name        string    `json:"name"`
	ImageURL    *string   `json:"imageUrl"`
	FriendCount int       `json:"friendCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UserGetResponse struct {
	ID          string    `json:"userId"`
	Name        string    `json:"name"`
	ImageURL    *string   `json:"imageUrl"`
	FriendCount int       `json:"friendCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UserCommentResponse struct {
	ID          *string `json:"userId"`
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageUrl"`
	FriendCount *int    `json:"friendCount"`
}
