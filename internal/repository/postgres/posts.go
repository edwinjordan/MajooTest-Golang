package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type PostsRepository struct {
	// Add necessary fields, e.g., DB connection
	Conn *pgxpool.Pool
}

func NewPostsRepository(conn *pgxpool.Pool) *PostsRepository {
	return &PostsRepository{Conn: conn}
}

func (r *PostsRepository) CreatePosts(ctx context.Context, post *domain.CreatePostsRequest) (*domain.Posts, error) {
	query := `
		INSERT INTO posts (title, content, slug, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	var id uuid.UUID
	err := r.Conn.QueryRow(ctx, query, post.Title, post.Content, utils.Slugify(post.Title)).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.Posts{
		ID:      id.String(),
		Title:   post.Title,
		Content: post.Content,
		Slug:    utils.Slugify(post.Title),
	}, nil
}

func (u *PostsRepository) GetPostsList(ctx context.Context, filter *domain.PostsFilter) ([]domain.Posts, error) {
	query := `
		SELECT
			u.id,
			u.title,
			u.content,
			u.slug,
            u.created_at,
            u.updated_at
		FROM posts u
        WHERE u.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(u.title ILIKE $1 OR u.content ILIKE $1)`)
		args = append(args, "%"+filter.Search+"%")
	}

	if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}
	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Posts

	for rows.Next() {
		var post domain.Posts
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Slug,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (u *PostsRepository) GetPosts(ctx context.Context, id uuid.UUID) (*domain.Posts, error) {
	tracer := otel.Tracer("repo.posts")
	ctx, span := tracer.Start(ctx, "PostsRepository.GetPosts")
	defer span.End()

	query := `
		SELECT
			id,
			title,
			content,
			slug,
			created_at,
			updated_at
		FROM posts
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := u.Conn.QueryRow(ctx, query, id)

	var post domain.Posts
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Slug,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		//	u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "error").Inc()
		return nil, err
	}

	//u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "success").Inc()
	return &post, nil
}

func (u *PostsRepository) UpdatePosts(ctx context.Context, id uuid.UUID, post *domain.Posts) (*domain.Posts, error) {
	query := `
		UPDATE posts
		SET title = $1,
			content = $2,
			slug = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, title, content, slug, created_at, updated_at`

	var updatedPost domain.Posts
	err := u.Conn.QueryRow(ctx, query, post.Title, post.Content, post.Slug, id).Scan(
		&updatedPost.ID,
		&updatedPost.Title,
		&updatedPost.Content,
		&updatedPost.Slug,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &updatedPost, nil
}

func (u *PostsRepository) DeletePosts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE posts
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := u.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
