package domain

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// CSVJobStatus represents the status of a CSV processing job
type CSVJobStatus string

const (
	CSVJobStatusPending    CSVJobStatus = "pending"
	CSVJobStatusProcessing CSVJobStatus = "processing"
	CSVJobStatusCompleted  CSVJobStatus = "completed"
	CSVJobStatusFailed     CSVJobStatus = "failed"
)

// CSVJob represents a CSV processing job
type CSVJob struct {
	ID            string       `json:"id" db:"id"`
	Filename      string       `json:"filename" db:"filename"`
	Status        CSVJobStatus `json:"status" db:"status"`
	TotalRows     int64        `json:"total_rows" db:"total_rows"`
	ProcessedRows int64        `json:"processed_rows" db:"processed_rows"`
	FailedRows    int64        `json:"failed_rows" db:"failed_rows"`
	ErrorMessage  *string      `json:"error_message,omitempty" db:"error_message"`
	StartedAt     time.Time    `json:"started_at" db:"started_at"`
	CompletedAt   time.Time    `json:"completed_at" db:"completed_at"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// CSVProcessingResult represents the result of processing a single CSV row
type CSVProcessingResult struct {
	RowNumber int                    `json:"row_number"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// CSVProcessingProgress represents real-time progress of CSV processing
type CSVProcessingProgress struct {
	JobID         string       `json:"job_id"`
	TotalRows     int64        `json:"total_rows"`
	ProcessedRows int64        `json:"processed_rows"`
	FailedRows    int64        `json:"failed_rows"`
	SuccessRate   float64      `json:"success_rate"`
	Status        CSVJobStatus `json:"status"`
	Message       string       `json:"message,omitempty"`
}

// CSVWorkerJob represents a job sent to CSV worker
type CSVWorkerJob struct {
	JobID     string
	RowNumber int
	Data      []string
	Headers   []string
}

// CSVWorkerResult represents result from CSV worker
type CSVWorkerResult struct {
	JobID     string
	RowNumber int
	Result    CSVProcessingResult
}

// CSVUploadRequest represents CSV upload request
type CSVUploadRequest struct {
	Files []*multipart.FileHeader `json:"files"`
}

// CSVUploadResponse represents CSV upload response
type CSVUploadResponse struct {
	Jobs    []CSVJob `json:"jobs"`
	Message string   `json:"message"`
}

// CSVRepository interface for CSV operations
type CSVRepository interface {
	CreateJob(ctx context.Context, job *CSVJob) error
	GetJobByID(ctx context.Context, id uuid.UUID) (*CSVJob, error)
	GetJobsByUserID(ctx context.Context) ([]*CSVJob, error)
	UpdateJobProgress(ctx context.Context, jobID string, processedRows, failedRows int64) error
	UpdateJobStatus(ctx context.Context, jobID string, status CSVJobStatus, errorMessage *string) error
	CompleteJob(ctx context.Context, jobID string, totalRows, processedRows, failedRows int64) error
}

// CSVService interface for CSV processing operations
type CSVService interface {
	UploadAndProcessCSV(ctx context.Context, files []*multipart.FileHeader) (*CSVUploadResponse, error)
	GetJobProgress(ctx context.Context, jobID uuid.UUID) (*CSVProcessingProgress, error)
	GetUserJobs(ctx context.Context) ([]*CSVJob, error)
	ProcessCSVFile(ctx context.Context, jobID string, reader io.Reader) error
}
