package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/edwinjordan/MajooTest-Golang/domain"
)

const (
	// DefaultWorkerPoolSize is the default number of workers
	DefaultWorkerPoolSize = 10
	// DefaultJobChannelSize is the default size of job channel
	DefaultJobChannelSize = 100
	// DefaultResultChannelSize is the default size of result channel
	DefaultResultChannelSize = 100
	// ProgressUpdateInterval defines how often to update progress
	ProgressUpdateInterval = time.Second * 2
	// MaxMemoryUsage defines maximum memory usage per CSV file (100MB)
	MaxMemoryUsage = 100 * 1024 * 1024
)

type csvService struct {
	csvRepo    domain.CSVRepository
	logger     *logrus.Logger
	workerPool int
}

// CSVWorkerPool manages concurrent CSV processing
type CSVWorkerPool struct {
	jobChan    chan domain.CSVWorkerJob
	resultChan chan domain.CSVWorkerResult
	doneChan   chan bool
	wg         *sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *logrus.Logger
}

// NewCSVService creates a new CSV service
func NewCSVService(csvRepo domain.CSVRepository, logger *logrus.Logger) *csvService {
	return &csvService{
		csvRepo:    csvRepo,
		logger:     logger,
		workerPool: DefaultWorkerPoolSize,
	}
}

// UploadAndProcessCSV handles multiple CSV file uploads and starts processing
func (s *csvService) UploadAndProcessCSV(ctx context.Context, files []*multipart.FileHeader) (*domain.CSVUploadResponse, error) {

	if len(files) == 0 {
		return nil, domain.ErrBadParamInput
	}

	var jobs []domain.CSVJob
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(files))

	// Process each file concurrently
	for _, fileHeader := range files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			// Validate file size to prevent memory issues
			if fh.Size > MaxMemoryUsage {
				errChan <- fmt.Errorf("file %s exceeds maximum size limit", fh.Filename)
				return
			}

			// Create job
			job := &domain.CSVJob{}
			job.ID = uuid.New().String()
			job.Filename = fh.Filename
			job.Status = domain.CSVJobStatusPending

			// Create job in database
			if err := s.csvRepo.CreateJob(ctx, job); err != nil {
				s.logger.WithError(err).Error("Failed to create CSV job")
				errChan <- err
				return
			}

			// Add to jobs list thread-safely
			mu.Lock()
			jobs = append(jobs, *job)
			mu.Unlock()

			// Start processing asynchronously
			go func(jobID string, fileHeader *multipart.FileHeader) {
				file, err := fileHeader.Open()
				if err != nil {
					s.logger.WithError(err).WithField("job_id", jobID).Error("Failed to open CSV file")
					errMsg := err.Error()
					s.csvRepo.UpdateJobStatus(context.Background(), jobID, domain.CSVJobStatusFailed, &errMsg)
					return
				}
				defer file.Close()

				if err := s.ProcessCSVFile(context.Background(), jobID, file); err != nil {
					s.logger.WithError(err).WithField("job_id", jobID).Error("CSV processing failed")
				}
			}(job.ID, fh)

		}(fileHeader)
	}

	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &domain.CSVUploadResponse{
		Jobs:    jobs,
		Message: fmt.Sprintf("Successfully uploaded %d CSV files for processing", len(jobs)),
	}, nil
}

// ProcessCSVFile processes a single CSV file using worker pool pattern
func (s *csvService) ProcessCSVFile(ctx context.Context, jobID string, reader io.Reader) error {
	// Update job status to processing
	if err := s.csvRepo.UpdateJobStatus(ctx, jobID, domain.CSVJobStatusProcessing, nil); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Create CSV reader
	csvReader := csv.NewReader(reader)
	csvReader.ReuseRecord = true // Memory optimization

	// Read headers
	headers, err := csvReader.Read()
	if err != nil {
		errMsg := "failed to read CSV headers"
		s.csvRepo.UpdateJobStatus(ctx, jobID, domain.CSVJobStatusFailed, &errMsg)
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Create worker pool
	pool := s.createWorkerPool(ctx, jobID)
	defer pool.Close()

	// Start result processor
	var processedRows, failedRows int64
	var totalRows int64

	resultProcessor := make(chan bool)
	go func() {
		defer close(resultProcessor)
		ticker := time.NewTicker(ProgressUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case result, ok := <-pool.resultChan:
				if !ok {
					return
				}

				if result.Result.Success {
					atomic.AddInt64(&processedRows, 1)
				} else {
					atomic.AddInt64(&failedRows, 1)
					s.logger.WithFields(logrus.Fields{
						"job_id":     jobID,
						"row_number": result.RowNumber,
						"error":      result.Result.Error,
					}).Warn("CSV row processing failed")
				}

			case <-ticker.C:
				// Periodic progress update
				current := atomic.LoadInt64(&processedRows)
				failed := atomic.LoadInt64(&failedRows)
				if err := s.csvRepo.UpdateJobProgress(ctx, jobID, current, failed); err != nil {
					s.logger.WithError(err).Error("Failed to update job progress")
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	// Read and send jobs to workers
	rowNumber := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.logger.WithError(err).WithField("job_id", jobID).Error("Error reading CSV row")
			atomic.AddInt64(&failedRows, 1)
			continue
		}

		rowNumber++
		totalRows++

		// Create a copy of the record to avoid race conditions
		recordCopy := make([]string, len(record))
		copy(recordCopy, record)

		job := domain.CSVWorkerJob{
			JobID:     jobID,
			RowNumber: rowNumber,
			Data:      recordCopy,
			Headers:   headers,
		}

		select {
		case pool.jobChan <- job:
			// Job sent successfully
		case <-ctx.Done():
			errMsg := "processing cancelled"
			s.csvRepo.UpdateJobStatus(ctx, jobID, domain.CSVJobStatusFailed, &errMsg)
			return ctx.Err()
		}
	}

	// Close job channel and wait for completion
	close(pool.jobChan)
	pool.wg.Wait()
	close(pool.resultChan)

	// Wait for result processor to finish
	<-resultProcessor

	// Final progress update and job completion
	finalProcessed := atomic.LoadInt64(&processedRows)
	finalFailed := atomic.LoadInt64(&failedRows)

	if err := s.csvRepo.CompleteJob(ctx, jobID, totalRows, finalProcessed, finalFailed); err != nil {
		return fmt.Errorf("failed to complete job: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"job_id":         jobID,
		"total_rows":     totalRows,
		"processed_rows": finalProcessed,
		"failed_rows":    finalFailed,
	}).Info("CSV processing completed")

	// Force garbage collection to free memory
	runtime.GC()

	return nil
}

// createWorkerPool creates and starts a worker pool
func (s *csvService) createWorkerPool(ctx context.Context, jobID string) *CSVWorkerPool {
	poolCtx, cancel := context.WithCancel(ctx)

	pool := &CSVWorkerPool{
		jobChan:    make(chan domain.CSVWorkerJob, DefaultJobChannelSize),
		resultChan: make(chan domain.CSVWorkerResult, DefaultResultChannelSize),
		doneChan:   make(chan bool),
		wg:         &sync.WaitGroup{},
		ctx:        poolCtx,
		cancel:     cancel,
		logger:     s.logger,
	}

	// Start workers
	for i := 0; i < s.workerPool; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker processes CSV jobs
func (wp *CSVWorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	wp.logger.WithField("worker_id", workerID).Debug("CSV worker started")
	defer wp.logger.WithField("worker_id", workerID).Debug("CSV worker stopped")

	for {
		select {
		case job, ok := <-wp.jobChan:
			if !ok {
				return // Channel closed
			}

			result := wp.processRow(job)

			select {
			case wp.resultChan <- domain.CSVWorkerResult{
				JobID:     job.JobID,
				RowNumber: job.RowNumber,
				Result:    result,
			}:
			case <-wp.ctx.Done():
				return
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

// processRow processes a single CSV row
func (wp *CSVWorkerPool) processRow(job domain.CSVWorkerJob) domain.CSVProcessingResult {
	// Simulate data processing and validation
	data := make(map[string]interface{})

	// Basic validation - ensure we have data for all headers
	if len(job.Data) != len(job.Headers) {
		return domain.CSVProcessingResult{
			RowNumber: job.RowNumber,
			Success:   false,
			Error:     fmt.Sprintf("column count mismatch: expected %d, got %d", len(job.Headers), len(job.Data)),
		}
	}

	// Map headers to data
	for i, header := range job.Headers {
		if i < len(job.Data) {
			data[header] = job.Data[i]
		}
	}

	// Simulate some processing time and potential errors
	if job.RowNumber%100 == 0 { // Simulate 1% error rate
		return domain.CSVProcessingResult{
			RowNumber: job.RowNumber,
			Success:   false,
			Error:     "simulated processing error",
		}
	}

	return domain.CSVProcessingResult{
		RowNumber: job.RowNumber,
		Success:   true,
		Data:      data,
	}
}

// Close closes the worker pool
func (wp *CSVWorkerPool) Close() {
	wp.cancel()
}

// GetJobProgress retrieves the progress of a CSV processing job
func (s *csvService) GetJobProgress(ctx context.Context, jobID uuid.UUID) (*domain.CSVProcessingProgress, error) {
	job, err := s.csvRepo.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	var successRate float64
	if job.TotalRows > 0 {
		successRate = float64(job.ProcessedRows-job.FailedRows) / float64(job.TotalRows) * 100
	}

	progress := &domain.CSVProcessingProgress{
		JobID:         job.ID,
		TotalRows:     job.TotalRows,
		ProcessedRows: job.ProcessedRows,
		FailedRows:    job.FailedRows,
		SuccessRate:   successRate,
		Status:        job.Status,
	}

	// Add status message
	switch job.Status {
	case domain.CSVJobStatusPending:
		progress.Message = "Job is waiting to be processed"
	case domain.CSVJobStatusProcessing:
		progress.Message = "Job is currently being processed"
	case domain.CSVJobStatusCompleted:
		progress.Message = "Job processing completed successfully"
	case domain.CSVJobStatusFailed:
		progress.Message = "Job processing failed"
		if job.ErrorMessage != nil {
			progress.Message = *job.ErrorMessage
		}
	}

	return progress, nil
}

// GetUserJobs retrieves all CSV jobs for a user
func (s *csvService) GetUserJobs(ctx context.Context) ([]*domain.CSVJob, error) {
	return s.csvRepo.GetJobsByUserID(ctx)
}
