package posts

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreatePostPayload struct {
	UserID     string
	PostInHTML string   `json:"postInHtml"`
	Tags       []string `json:"tags"`
}

func (p CreatePostPayload) Validate() error {
	for i := range p.Tags {
		if p.Tags[i] == "" {
			return errors.New("tags must not be empty")
		}
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
		validation.Field(&p.PostInHTML, validation.Required),
		validation.Field(&p.Tags, validation.Required),
	)
}