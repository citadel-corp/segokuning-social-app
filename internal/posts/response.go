package posts

import (
	"time"

	"github.com/citadel-corp/segokuning-social-app/internal/comments"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
	"github.com/citadel-corp/segokuning-social-app/internal/user"
)

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
	Error   string
}

var (
	SuccessCreateResponse        = Response{Code: 200, Message: "Post created successfully"}
	SuccessCreateCommentResponse = Response{Code: 200, Message: "Comment created successfully"}
	SuccessListResponse          = Response{Code: 200, Message: "Posts fetched successfully"}
)

type ListPostResponse struct {
	PostID   string                     `json:"postId"`
	Post     PostResponse               `json:"post"`
	Comments []comments.CommentResponse `json:"comments"`
	User     user.UserGetResponse       `json:"creator"`
}

type PostResponse struct {
	ID        string    `json:"-"`
	Content   string    `json:"postInHtml"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}
