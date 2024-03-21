package posts

import "github.com/citadel-corp/segokuning-social-app/internal/common/response"

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
	Error   error
}

var (
	SuccessCreateResponse        = Response{Code: 200, Message: "Post created successfully"}
	SuccessCreateCommentResponse = Response{Code: 200, Message: "Comment created successfully"}
)
