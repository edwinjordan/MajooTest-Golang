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

type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{Conn: conn}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error) {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, name, email, created_at, updated_at`

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	var createdUser domain.User
	err = u.Conn.QueryRow(ctx, query, user.Name, user.Email, hashedPassword).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (u *UserRepository) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
            u.created_at,
            u.updated_at
		FROM users u
        WHERE u.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(u.name ILIKE $1 OR u.email ILIKE $1)`)
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

	var users []domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.GetUser")
	defer span.End()

	query := `
		SELECT
			id,
			name,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := u.Conn.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		//	u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "error").Inc()
		return nil, err
	}

	//u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "success").Inc()
	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users
		SET name = $1,
			email = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, name, email, created_at, updated_at`

	var updatedUser domain.User
	err := u.Conn.QueryRow(ctx, query, user.Name, user.Email, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
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
