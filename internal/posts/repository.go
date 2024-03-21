package posts

import (
	"context"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	"github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Repository interface {
	Create(ctx context.Context, post *Posts) error
	GetByID(ctx context.Context, id string) (*Posts, error)
	CreateComment(ctx context.Context, comment *Comment) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, post *Posts) error {
	id, err := gonanoid.Generate("abcdef1234567890", 16)
	if err != nil {
		return err
	}

	_, err = d.db.DB().ExecContext(ctx, `
			INSERT INTO posts (
				id, user_id, content, tags
			) VALUES (
				$1, $2, $3, $4
			)
		`, id, post.UserID, post.Content, post.Tags)
	if err != nil {
		return err
	}

	return nil
}

// GetByID implements Repository.
func (d *dbRepository) GetByID(ctx context.Context, id string) (*Posts, error) {
	row := d.db.DB().QueryRowContext(ctx, `
		SELECT id, user_id, content, tags, created_at
		FROM posts
		WHERE id = $1;
	`, id)

	p := &Posts{}
	err := row.Scan(&p.ID, &p.UserID, &p.Content, pq.Array(&p.Tags), &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *dbRepository) CreateComment(ctx context.Context, comment *Comment) error {
	row := d.db.DB().QueryRowContext(ctx, `
			INSERT INTO comments (
				user_id, post_id, content
			) VALUES (
				$1, $2, $3
			)
			RETURNING id
		`, comment.UserID, comment.PostID, comment.Content)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}
