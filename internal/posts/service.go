package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type Service interface {
	Create(ctx context.Context, req CreatePostPayload) Response
	List(ctx context.Context, req ListPostPayload) Response
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

func (s *postsService) List(ctx context.Context, req ListPostPayload) Response {
	serviceName := "posts.List"
	posts, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error(fmt.Sprintf("[%s] error while fetching posts: %s", serviceName, err.Error()))
			return ErrorInternal
		}
	}

	resp := SuccessListResponse
	resp.Data = posts
	resp.Meta = pagination

	return resp
}
