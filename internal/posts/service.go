package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	userfriends "github.com/citadel-corp/segokuning-social-app/internal/user_friends"
)

type Service interface {
	Create(ctx context.Context, req CreatePostPayload) Response
	CreatePostComment(ctx context.Context, req CreatePostCommentPayload) Response
}

type postsService struct {
	repository            Repository
	userFriendsRepository userfriends.Repository
}

func NewService(repository Repository, userFriendsRepository userfriends.Repository) Service {
	return &postsService{repository: repository, userFriendsRepository: userFriendsRepository}
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

func (s *postsService) CreatePostComment(ctx context.Context, req CreatePostCommentPayload) Response {
	serviceName := "posts.CreatePostComment"

	//validate post is not found
	post, err := s.repository.GetByID(ctx, req.PostID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorNotFound
	}
	if err != nil {
		return ErrorInternal
	}
	//validate post creator is users friend
	_, err = s.userFriendsRepository.GetByFriendID(ctx, req.UserID, post.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorBadRequest
	}
	if err != nil {
		return ErrorInternal
	}

	comment := &Comment{
		UserID:  req.UserID,
		PostID:  req.PostID,
		Content: req.Comment,
	}

	err = s.repository.CreateComment(ctx, comment)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] error creating post: %v", serviceName, err))
		return ErrorInternal
	}

	return SuccessCreateCommentResponse
}
