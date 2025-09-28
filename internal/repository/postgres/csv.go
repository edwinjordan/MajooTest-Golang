package postgres

import (
	"context"
	"time"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type csvRepository struct {
	Conn *pgxpool.Pool
}

// NewCSVRepository creates a new CSV repository
func NewCSVRepository(conn *pgxpool.Pool) *csvRepository {
	return &csvRepository{Conn: conn}
}

// CreateJob creates a new CSV processing job
func (r *csvRepository) CreateJob(ctx context.Context, job *domain.CSVJob) error {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now

	query := `
		INSERT INTO csv_jobs ( filename, status, total_rows, processed_rows, failed_rows, error_message, started_at, completed_at, created_at, updated_at)
		VALUES ( $1, $2, $3, $4, $5, $6, Now(), Now(), Now(), Now())
		RETURNING id`
	var id uuid.UUID
	err := r.Conn.QueryRow(ctx, query, job.Filename, job.Status, job.TotalRows, job.ProcessedRows, job.FailedRows, job.ErrorMessage).Scan(&id)

	if err != nil {
		return err
	}

	// populate the passed job with the returned id and timestamps
	job.ID = id.String()
	job.StartedAt = now
	job.CompletedAt = now

	return nil
}

// GetJobByID retrieves a CSV job by ID
func (r *csvRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*domain.CSVJob, error) {
	tracer := otel.Tracer("repo.csv")
	ctx, span := tracer.Start(ctx, "CSVRepository.GetJobByID")
	defer span.End()

	query := `
		SELECT
			id,
			filename, 
			status, 
			total_rows, 
			processed_rows, 
			failed_rows, 
			error_message, 
			started_at, 
			completed_at, 
			created_at, 
			updated_at
		FROM csv_jobs
		WHERE id = $1`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := r.Conn.QueryRow(ctx, query, id)

	var job domain.CSVJob
	err := row.Scan(
		&job.ID,
		&job.Filename,
		&job.Status,
		&job.TotalRows,
		&job.ProcessedRows,
		&job.FailedRows,
		&job.ErrorMessage,
		&job.StartedAt,
		&job.CompletedAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		//	u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "error").Inc()
		return nil, err
	}

	//u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "success").Inc()
	return &job, nil
}

// GetJobsByUserID retrieves all CSV jobs for a user
// GetJobsByUserID retrieves all CSV jobs for a user
func (r *csvRepository) GetJobsByUserID(ctx context.Context) ([]*domain.CSVJob, error) {
	tracer := otel.Tracer("repo.csv")
	ctx, span := tracer.Start(ctx, "CSVRepository.GetJobsByUserID")
	defer span.End()

	query := `
		SELECT id,  filename, status, total_rows, processed_rows, failed_rows, error_message, started_at, completed_at, created_at, updated_at
		FROM csv_jobs
	
		ORDER BY created_at DESC
	`

	span.SetAttributes(attribute.String("query.statement", query))
	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.CSVJob
	for rows.Next() {
		var job domain.CSVJob
		err := rows.Scan(
			&job.ID,
			&job.Filename,
			&job.Status,
			&job.TotalRows,
			&job.ProcessedRows,
			&job.FailedRows,
			&job.ErrorMessage,
			&job.StartedAt,
			&job.CompletedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}
	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, err
	}
	return jobs, nil
}

// UpdateJobProgress updates the progress of a CSV processing job
func (r *csvRepository) UpdateJobProgress(ctx context.Context, jobID string, processedRows, failedRows int64) error {
	query := `
		UPDATE csv_jobs
		SET processed_rows = $2, failed_rows = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.Conn.Exec(ctx, query, jobID, processedRows, failedRows, time.Now())
	return err
}

// UpdateJobStatus updates the status of a CSV processing job
func (r *csvRepository) UpdateJobStatus(ctx context.Context, jobID string, status domain.CSVJobStatus, errorMessage *string) error {
	now := time.Now()

	var startedAt *time.Time
	var completedAt *time.Time

	switch status {
	case domain.CSVJobStatusProcessing:
		startedAt = &now
	case domain.CSVJobStatusCompleted, domain.CSVJobStatusFailed:
		completedAt = &now
	}

	query := `
		UPDATE csv_jobs
		SET status = $2, error_message = $3, started_at = COALESCE($4, started_at), completed_at = COALESCE($5, completed_at), updated_at = $6
		WHERE id = $1
	`

	_, err := r.Conn.Exec(ctx, query, jobID, status, errorMessage, startedAt, completedAt, now)
	return err
}

// CompleteJob marks a CSV processing job as completed
func (r *csvRepository) CompleteJob(ctx context.Context, jobID string, totalRows, processedRows, failedRows int64) error {
	now := time.Now()

	query := `
		UPDATE csv_jobs
		SET status = $2, total_rows = $3, processed_rows = $4, failed_rows = $5, completed_at = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.Conn.Exec(ctx, query, jobID, domain.CSVJobStatusCompleted, totalRows, processedRows, failedRows, now, now)
	return err
}
