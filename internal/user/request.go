package user

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var phoneNumberValidationRule = validation.NewStringRule(func(s string) bool {
	return strings.HasPrefix(s, "+")
}, "phone number must start with international calling code")

type CreateUserPayload struct {
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
	Name            string `json:"name"`
	Password        string `json:"password"`
}

func (p CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CredentialType, validation.Required, validation.In("phone", "email")),
		validation.Field(&p.CredentialValue, validation.Required, validation.
			When(p.CredentialType == "email", is.EmailFormat).
			Else(phoneNumberValidationRule, validation.Length(7, 13))),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}

type LoginPayload struct {
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
	Password        string `json:"password"`
}

func (p LoginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CredentialType, validation.Required, validation.In("phone", "email")),
		validation.Field(&p.CredentialValue, validation.Required, validation.
			When(p.CredentialType == "email", is.EmailFormat).
			Else(phoneNumberValidationRule, validation.Length(7, 13))),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}

type LinkEmailPayload struct {
	Email string `json:"email"`
}

func (p LinkEmailPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.EmailFormat),
	)
}

type LinkPhoneNumberPayload struct {
	Phone string `json:"phone"`
}

func (p LinkPhoneNumberPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Phone, validation.Required, phoneNumberValidationRule, validation.Length(7, 13)),
	)
}

type UpdateUserPayload struct {
	ImageURL string `json:"imageUrl"`
	Name     string `json:"name"`
}

func (p UpdateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ImageURL, validation.Required, is.URL),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
	)
}

type sortBy string
type userSortBy sortBy

var (
	SortByFriendCount userSortBy = "friendCount"
	SortByCreatedAt   userSortBy = "createdAt"
)

var userSortBys []interface{} = []interface{}{SortByFriendCount, SortByCreatedAt}

type ListUserPayload struct {
	OnlyFriend bool `schema:"onlyFriend" binding:"omitempty"`
	UserID     string
	Search     string     `schema:"search" binding:"omitempty"`
	Limit      int        `schema:"limit" binding:"omitempty"`
	Offset     int        `schema:"offset" binding:"omitempty"`
	SortBy     userSortBy `schema:"sortBy" binding:"omitempty"`
	OrderBy    string     `schema:"orderBy" binding:"omitempty"`
}

func (p ListUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.SortBy, validation.In(userSortBys...)),
		validation.Field(&p.OrderBy, validation.In("asc", "desc")),
		validation.Field(&p.Limit, validation.When(p.Offset != 0, validation.Required)),
		validation.Field(&p.Offset, validation.When(p.Limit != 0, validation.NotNil)),
	)
}
