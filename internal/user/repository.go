package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateWithEmail(ctx context.Context, user *User) error
	CreateWithPhoneNumber(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	List(ctx context.Context, filter ListUserPayload) ([]UserListResponse, *response.Pagination, error)
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

func (d *dbRepository) List(ctx context.Context, filter ListUserPayload) ([]UserListResponse, *response.Pagination, error) {
	var users []UserListResponse
	var pagination *response.Pagination

	var (
		selectStatement     string
		whereStatement      string
		query               string
		joinStatement       string
		orderStatement      string
		paginationStatement string
		args                []interface{}
		columnCtr           int = 1
	)

	if filter.OnlyFriend && filter.UserID != "" {
		whereStatement = fmt.Sprintf("%s WHERE user_friends.friend_id = $%d", whereStatement, columnCtr)
		joinStatement = fmt.Sprintf("%s JOIN user_friends ON users.id = user_friends.user_id", joinStatement)
		args = append(args, filter.UserID)
		columnCtr++
	}

	if filter.Search != "" {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s lower(users.name) LIKE CONCAT('%%',$%d::text,'%%')", whereStatement, columnCtr)
		args = append(args, strings.ToLower(filter.Search))
		columnCtr++
	}

	var orderBy string
	switch filter.OrderBy {
	case "asc":
		orderBy = "asc"
	case "desc":
		orderBy = "desc"
	default:
		orderBy = "desc"
	}

	switch filter.SortBy {
	case userSortBy(SortByFriendCount):
		orderStatement = fmt.Sprintf("%s ORDER BY users.friend_count %s", orderStatement, orderBy)
	case userSortBy(SortByCreatedAt):
		orderStatement = fmt.Sprintf("%s ORDER BY users.created_at %s", orderStatement, orderBy)
	default:
		orderStatement = fmt.Sprintf("%s ORDER BY users.created_at %s", orderStatement, orderBy)
	}

	if filter.Limit == 0 {
		filter.Limit = 5
	}

	var rows *sql.Rows
	var err error
	pagination = &response.Pagination{
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	selectStatement = fmt.Sprintf(`
		SELECT COUNT(*) OVER() AS total_count, users.id as productId, users.name as name, users.image_url as imageUrl,
			users.friend_count as friendCount, users.created_at as createdAt
		FROM users
	%s`, selectStatement)

	paginationStatement = fmt.Sprintf("%s LIMIT $%d", paginationStatement, columnCtr)
	args = append(args, filter.Limit)
	columnCtr++

	paginationStatement = fmt.Sprintf("%s OFFSET $%d", paginationStatement, columnCtr)
	args = append(args, filter.Offset)

	query = fmt.Sprintf("%s %s %s %s %s;", selectStatement, joinStatement, whereStatement, orderStatement, paginationStatement)

	rows, err = d.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var u UserListResponse
		if err := rows.Scan(&pagination.Total, &u.ID, &u.Name, &u.ImageURL, &u.FriendCount, &u.CreatedAt); err != nil {
			return users, nil, err
		}
		users = append(users, u)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return users, nil, err
	}

	return users, pagination, nil
}

func insertWhereStatement(condition bool, statement string) string {
	if condition {
		return fmt.Sprintf(`%v AND`, statement)
	}
	return fmt.Sprintf(`%v WHERE`, statement)
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
