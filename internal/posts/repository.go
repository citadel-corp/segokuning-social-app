package posts

import (
	"context"

	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Repository interface {
	Create(ctx context.Context, post *Posts) error
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
