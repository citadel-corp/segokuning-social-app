package userfriends

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	Create(ctx context.Context, req CreateUserFriendPayload) Response
}

type userFriendsService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userFriendsService{repository: repository}
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

	var pgErr *pgconn.PgError
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] error while adding friend: %s", serviceName, err.Error()))
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrFriendAlreadyExists
			case "23503":
				return ErrFriendNotExists
			default:
				return ErrorInternal
			}
		}

		return ErrorInternal
	}

	return SuccessCreateResponse
}
