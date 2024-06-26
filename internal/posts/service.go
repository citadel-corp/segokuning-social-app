package posts

import (
	"context"
	"database/sql"
	"errors"

	userfriends "github.com/citadel-corp/segokuning-social-app/internal/user_friends"
)

type Service interface {
	Create(ctx context.Context, req CreatePostPayload) Response
	CreatePostComment(ctx context.Context, req CreatePostCommentPayload) Response
	List(ctx context.Context, req ListPostPayload) Response
}

type postsService struct {
	repository            Repository
	userFriendsRepository userfriends.Repository
}

func NewService(repository Repository, userFriendsRepository userfriends.Repository) Service {
	return &postsService{repository: repository, userFriendsRepository: userFriendsRepository}
}

func (s *postsService) Create(ctx context.Context, req CreatePostPayload) Response {
	var resp Response

	post := &Posts{
		UserID:  req.UserID,
		Content: req.PostInHTML,
		Tags:    req.Tags,
	}

	err := s.repository.Create(ctx, post)
	if err != nil {
		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessCreateResponse
}

func (s *postsService) CreatePostComment(ctx context.Context, req CreatePostCommentPayload) Response {
	var resp Response

	//validate post is not found
	post, err := s.repository.GetByID(ctx, req.PostID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorNotFound
	}
	if err != nil {
		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	//validate post creator is users friend
	if post.UserID != req.UserID {
		_, err = s.userFriendsRepository.GetByFriendID(ctx, req.UserID, post.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorBadRequest
		}
		if err != nil {
			resp = ErrorInternal
			resp.Error = err.Error()
			return resp
		}
	}

	comment := &Comment{
		UserID:  req.UserID,
		PostID:  req.PostID,
		Content: req.Comment,
	}

	err = s.repository.CreateComment(ctx, comment)
	if err != nil {
		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessCreateCommentResponse
}

func (s *postsService) List(ctx context.Context, req ListPostPayload) Response {
	var resp Response

	posts, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			resp = ErrorInternal
			resp.Error = err.Error()
			return resp
		}
	}

	resp = SuccessListResponse
	resp.Data = posts
	resp.Meta = pagination

	return resp
}
