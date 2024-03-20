package userfriends

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	serviceName := "userfriends.Create"

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
				slog.Error(fmt.Sprintf("[%s] error while adding friend: %s", serviceName, err.Error()))
				return ErrorInternal
			}
		}

		slog.Error(fmt.Sprintf("[%s] error while adding friend: %s", serviceName, err.Error()))
		return ErrorInternal
	}

	return SuccessCreateResponse
}

func (s *userFriendsService) Delete(ctx context.Context, req DeleteUserFriendPayload) Response {
	serviceName := "userfriends.Delete"

	if req.UserID == req.LoggedUserID {
		return ErrCannotDeleteSelf
	}

	// get user id (friend id)
	friend, err := s.userRepository.GetByID(ctx, req.UserID)
	if friend == nil {
		return ErrFriendNotExists
	}
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] error while getting friend detail: %s", serviceName, err.Error()))
		return ErrorInternal
	}

	// check friendship
	friendship, err := s.repository.GetByFriendID(ctx, req.LoggedUserID, req.UserID)
	if friendship == nil {
		return ErrNotFriend
	}
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] error while checking friendship: %s", serviceName, err.Error()))
		return ErrorInternal
	}

	err = s.repository.RemoveFriend(ctx, req.LoggedUserID, req.UserID)
	if err != nil {
		return ErrorInternal
	}

	return SuccessDeleteResponse
}
