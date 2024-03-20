package userfriends

import (
	"context"
	"database/sql"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
)

type Repository interface {
	AddFriend(ctx context.Context, userFriend *UserFriends) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) AddFriend(ctx context.Context, userFriend *UserFriends) error {
	err := d.db.StartTx(ctx, func(tx *sql.Tx) error {
		// add user's friend
		_, err := d.db.DB().ExecContext(ctx, `
				INSERT INTO user_friends (
					user_id, friend_id
				) VALUES (
					$1, $2
				)
			`, userFriend.UserID, userFriend.FriendID)
		if err != nil {
			return err
		}

		// add user as friend's friend
		_, err = d.db.DB().ExecContext(ctx, `
				INSERT INTO user_friends (
					user_id, friend_id
				) VALUES (
					$1, $2
				)
			`, userFriend.FriendID, userFriend.UserID)
		if err != nil {
			return err
		}

		// add friendCount to user and friend
		_, err = d.db.DB().ExecContext(ctx, `
				UPDATE users
				SET friend_count = friend_count + 1
				WHERE id = $1
				OR id = $2
			`, userFriend.UserID, userFriend.FriendID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
