package userfriends

import "github.com/citadel-corp/segokuning-social-app/internal/common/response"

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
	Error   string
}

var (
	SuccessCreateResponse = Response{Code: 200, Message: "Friend added successfully"}
	SuccessDeleteResponse = Response{Code: 200, Message: "Friend deleted successfully"}
)
