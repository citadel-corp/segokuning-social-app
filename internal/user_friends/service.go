package userfriends

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/segokuning-social-app/internal/user"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	Create(ctx context.Context, req CreateUserFriendPayload) Response
	Delete(ctx context.Context, req DeleteUserFriendPayload) Response
}

type userFriendsService struct {
	repository     Repository
	userRepository user.Repository
}

func NewService(repository Repository, userRepository user.Repository) Service {
	return &userFriendsService{repository: repository, userRepository: userRepository}
}

func (s *userFriendsService) Create(ctx context.Context, req CreateUserFriendPayload) Response {
	var resp Response

	userFriend := &UserFriends{
		UserID:   req.LoggedUserID,
		FriendID: req.UserID,
	}

	if userFriend.UserID == userFriend.FriendID {
		return ErrCannotAddSelf
	}

	err := s.repository.AddFriend(ctx, userFriend)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrFriendAlreadyExists
			case "23503":
				return ErrFriendNotExists
			default:
				resp = ErrorInternal
				resp.Error = err.Error()
				return resp
			}
		}

		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessCreateResponse
}

func (s *userFriendsService) Delete(ctx context.Context, req DeleteUserFriendPayload) Response {
	var resp Response

	if req.UserID == req.LoggedUserID {
		return ErrCannotDeleteSelf
	}

	// get user id (friend id)
	_, err := s.userRepository.GetByID(ctx, req.UserID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return ErrFriendNotExists
		}

		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	// check friendship
	_, err = s.repository.GetByFriendID(ctx, req.LoggedUserID, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFriend
		}

		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	err = s.repository.RemoveFriend(ctx, req.LoggedUserID, req.UserID)
	if err != nil {
		resp = ErrorInternal
		resp.Error = err.Error()
		return resp
	}

	return SuccessDeleteResponse
}
