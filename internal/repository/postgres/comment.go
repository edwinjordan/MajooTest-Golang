package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type CommentRepository struct {
	// Add necessary fields, e.g., DB connection
	Conn *pgxpool.Pool
}

func NewCommentRepository(conn *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{Conn: conn}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *domain.CreateCommentRequest) (*domain.Comment, error) {
	query := `
		INSERT INTO comments (post_id, user_id, body, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	var id uuid.UUID
	err := r.Conn.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Body).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.Comment{
		ID:        id.String(),
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Body:      comment.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *CommentRepository) GetCommentList(ctx context.Context, filter *domain.CommentFilter) ([]domain.Comment, error) {
	query := `
		SELECT
			u.id,
			u.post_id,
			u.user_id,
			u.body,
			u.created_at,
			u.updated_at
		FROM comments u
        WHERE u.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(u.body ILIKE $1)`)
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

	var comments []domain.Comment

	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Body,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (u *CommentRepository) GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	tracer := otel.Tracer("repo.comments")
	ctx, span := tracer.Start(ctx, "CommentRepository.GetComment")
	defer span.End()

	query := `
		SELECT
			id,
			post_id,
			user_id,
			body,
			created_at,
			updated_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := u.Conn.QueryRow(ctx, query, id)

	var comment domain.Comment
	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		//	u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "error").Inc()
		return nil, err
	}

	//u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "success").Inc()
	return &comment, nil
}

func (u *CommentRepository) UpdateComment(ctx context.Context, id uuid.UUID, comment *domain.Comment) (*domain.Comment, error) {
	query := `
		UPDATE comments
		SET body = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, post_id, user_id, body, created_at, updated_at`

	var updatedComment domain.Comment
	err := u.Conn.QueryRow(ctx, query, comment.Body, id).Scan(
		&updatedComment.ID,
		&updatedComment.PostID,
		&updatedComment.UserID,
		&updatedComment.Body,
		&updatedComment.CreatedAt,
		&updatedComment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &updatedComment, nil
}

func (u *CommentRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments
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
