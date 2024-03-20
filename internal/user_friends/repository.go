package userfriends

import (
	"context"
	"database/sql"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
)

type Repository interface {
	AddFriend(ctx context.Context, userFriend *UserFriends) error
	RemoveFriend(ctx context.Context, userID string, friendID string) error
	GetByFriendID(ctx context.Context, userID string, friendID string) (*UserFriends, error)
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

func (d *dbRepository) RemoveFriend(ctx context.Context, userID string, friendID string) error {
	err := d.db.StartTx(ctx, func(tx *sql.Tx) error {
		_, err := d.db.DB().ExecContext(ctx, `
				DELETE FROM user_friends
				WHERE (user_id = $1 AND friend_id = $2)
				OR (user_id = $2 AND friend_id = $1)
			`, userID, friendID)
		if err != nil {
			return err
		}

		// add friendCount to user and friend
		_, err = d.db.DB().ExecContext(ctx, `
				UPDATE users
				SET friend_count = friend_count - 1
				WHERE id = $1
				OR id = $2
			`, userID, friendID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (d *dbRepository) GetByFriendID(ctx context.Context, userID string, friendID string) (*UserFriends, error) {
	getQuery := `
		SELECT id, user_id, friend_id, created_at FROM user_friends
		WHERE user_id = $1 AND friend_id = $2;
	`
	row := d.db.DB().QueryRowContext(ctx, getQuery, userID, friendID)

	var u UserFriends
	err := row.Scan(&u.ID, &u.UserID, &u.FriendID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
