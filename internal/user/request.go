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
			Else(validation.Length(7, 13), phoneNumberValidationRule)),
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
			Else(validation.Length(7, 13), phoneNumberValidationRule)),
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
		validation.Field(&p.Phone, validation.Required, validation.Length(7, 13), phoneNumberValidationRule),
	)
}
