package userfriends

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateUserFriendPayload struct {
	LoggedUserID string
	UserID       string `json:"userId"`
}

func (p CreateUserFriendPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.LoggedUserID, validation.Required.Error(ErrorUnauthorized.Message)),
		validation.Field(&p.UserID, validation.Required),
	)
}
