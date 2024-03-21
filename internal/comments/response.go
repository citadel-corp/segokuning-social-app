package comments

import (
	"time"

	"github.com/citadel-corp/segokuning-social-app/internal/user"
)

type CommentResponse struct {
	ID        *string                  `json:"-"`
	Content   *string                  `json:"comment"`
	User      user.UserCommentResponse `json:"creator"`
	CreatedAt *time.Time               `json:"createdAt"`
}
