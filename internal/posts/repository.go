package posts

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/citadel-corp/segokuning-social-app/internal/comments"
	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
	"github.com/citadel-corp/segokuning-social-app/internal/user"
	"github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Repository interface {
	Create(ctx context.Context, post *Posts) error
	GetByID(ctx context.Context, id string) (*Posts, error)
	CreateComment(ctx context.Context, comment *Comment) error
	List(ctx context.Context, filter ListPostPayload) ([]ListPostResponse, *response.Pagination, error)
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

func (d *dbRepository) List(ctx context.Context, filter ListPostPayload) ([]ListPostResponse, *response.Pagination, error) {
	var resp []ListPostResponse
	var pagination *response.Pagination

	var (
		withStatement   string
		selectStatement string
		query           string
		args            []interface{}
		columnCtr       int = 1
	)

	if filter.Limit == 0 {
		filter.Limit = 5
	}

	var rows *sql.Rows
	var err error
	pagination = &response.Pagination{
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	withStatement = fmt.Sprintf(`
		WITH p AS (
			SELECT COUNT(*) OVER() AS total_count, posts.*
			FROM posts
			LEFT JOIN user_friends uf ON uf.user_id = $%d
			AND posts.user_id = uf.friend_id
			WHERE (posts.user_id = $%d OR posts.user_id = uf.friend_id)
	`, columnCtr, columnCtr+1)
	args = append(args, filter.UserID)
	columnCtr++
	args = append(args, filter.UserID)
	columnCtr++

	if filter.Search != "" {
		withStatement = fmt.Sprintf("%s AND lower(posts.content) LIKE CONCAT('%%',$%d::text,'%%')", withStatement, columnCtr)
		args = append(args, strings.ToLower(filter.Search))
		columnCtr++
	}

	if len(filter.SearchTags) > 0 {
		for i := range filter.SearchTags {
			withStatement = fmt.Sprintf("%s AND $%d = ANY(posts.tags)", withStatement, columnCtr)
			args = append(args, filter.SearchTags[i])
			columnCtr++
		}
	}

	withStatement = fmt.Sprintf("%s LIMIT $%d OFFSET $%d) ", withStatement, columnCtr, columnCtr+1)

	args = append(args, filter.Limit)
	columnCtr++
	args = append(args, filter.Offset)
	columnCtr++

	selectStatement = `
		SELECT p.total_count, p.id as postId, p."content" as postInHtml, p.tags, p.created_at as product_created_at,
			c.id, c."content" as "comment", c.created_at as comment_created_at,
			pu.id as userId, pu.name as name, pu.image_url as imageUrl, pu.friend_count as friendCount,
			pu.created_at as user_created_at,
			cu.id as userId, cu.name as name, cu.image_url as imageUrl, cu.friend_count as friendCount
		FROM p
		JOIN users pu ON pu.id = p.user_id
		LEFT JOIN "comments" c ON p.id = c.post_id
		LEFT JOIN users cu ON cu.id = c.user_id 
		ORDER BY p.created_at desc, c.created_at desc
	`

	query = fmt.Sprintf("%s %s;", withStatement, selectStatement)

	rows, err = d.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	var ctrIndex int = 0
	resp = append(resp, ListPostResponse{})
	for rows.Next() {
		var p PostResponse
		var c comments.CommentResponse
		var pu user.UserGetResponse
		var cu user.UserCommentResponse
		if err := rows.Scan(&pagination.Total, &p.ID, &p.Content, pq.Array(&p.Tags), &p.CreatedAt,
			&c.ID, &c.Content, &c.CreatedAt,
			&pu.ID, &pu.Name, &pu.ImageURL, &pu.FriendCount, &pu.CreatedAt,
			&cu.ID, &cu.Name, &cu.ImageURL, &cu.FriendCount); err != nil {
			return resp, nil, err
		}

		if resp[ctrIndex].PostID == "" {
			resp[ctrIndex].PostID = p.ID
			resp[ctrIndex].Post = p
			resp[ctrIndex].User = pu
			resp[ctrIndex].Comments = []comments.CommentResponse{}
		} else if resp[ctrIndex].PostID != p.ID {
			ctrIndex++
			resp = append(resp, ListPostResponse{
				PostID:   p.ID,
				Post:     p,
				User:     pu,
				Comments: []comments.CommentResponse{},
			})
		}

		if c.ID != nil {
			c.User = cu
			resp[ctrIndex].Comments = append(resp[ctrIndex].Comments, c)
		}
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return resp, nil, err
	}

	return resp, pagination, nil
}
