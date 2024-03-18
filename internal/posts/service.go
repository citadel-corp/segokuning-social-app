package posts

import (
	"context"
	"fmt"
	"log/slog"
)

type Service interface {
	Create(ctx context.Context, req CreatePostPayload) Response
}

type postsService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &postsService{repository: repository}
}

func (s *postsService) Create(ctx context.Context, req CreatePostPayload) Response {
	serviceName := "posts.Create"

	post := &Posts{
		UserID:  req.UserID,
		Content: req.PostInHTML,
		Tags:    req.Tags,
	}

	err := s.repository.Create(ctx, post)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] error creating post: %s", serviceName, err.Error()))
		return ErrorInternal
	}

	return SuccessCreateResponse
}
