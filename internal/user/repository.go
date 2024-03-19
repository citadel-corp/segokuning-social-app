package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateWithEmail(ctx context.Context, user *User) error
	CreateWithPhoneNumber(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) CreateWithEmail(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			id, name, email, hashed_password
		) VALUES (
			$1, $2, $3, $4
		)
	`
	_, err := d.db.DB().ExecContext(ctx, createUserQuery, user.ID, user.Name, user.Email, user.HashedPassword)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUserEmailAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

// CreateWithPhoneNumber implements Repository.
func (d *dbRepository) CreateWithPhoneNumber(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			id, name, phone_number, hashed_password
		) VALUES (
			$1, $2, $3, $4
		)
	`
	_, err := d.db.DB().ExecContext(ctx, createUserQuery, user.ID, user.Name, user.PhoneNumber, user.HashedPassword)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUserPhoneNumberAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

// GetByEmail implements Repository.
func (d *dbRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	getUserQuery := `
		SELECT id, name, email, phone_number, friend_count, image_url, hashed_password, created_at FROM users
		WHERE email = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, email)
	return d.scanUser(row)
}

// GetByPhoneNumber implements Repository.
func (d *dbRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	getUserQuery := `
		SELECT id, name, email, phone_number, friend_count, image_url, hashed_password, created_at FROM users
		WHERE phone_number = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, phoneNumber)
	return d.scanUser(row)
}

func (d *dbRepository) GetByID(ctx context.Context, id string) (*User, error) {
	getUserQuery := `
		SELECT id, name, email, phone_number, friend_count, image_url, hashed_password, created_at FROM users
		WHERE id = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id)
	return d.scanUser(row)
}

// Update implements Repository.
func (d *dbRepository) Update(ctx context.Context, user *User) error {
	updateQuery := `
		UPDATE users
		SET name = $1,
		email = $2,
		phone_number = $3,
		friend_count = $4,
		image_url = $5
		WHERE id = $6;
	`
	_, err := d.db.DB().ExecContext(ctx, updateQuery, user.Name, user.Email, user.PhoneNumber, user.FriendCount, user.ImageURL, user.ID)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				if strings.Contains(pgErr.ConstraintName, "phone_number") {
					return ErrUserPhoneNumberAlreadyExists
				}
				if strings.Contains(pgErr.ConstraintName, "email") {
					return ErrUserEmailAlreadyExists
				}
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (d *dbRepository) scanUser(row *sql.Row) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PhoneNumber, &u.FriendCount, &u.ImageURL, &u.HashedPassword, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
