package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/edwinjordan/MajooTest-Golang/domain"
)

// CSVHandler handles CSV-related HTTP requests
type CSVHandler struct {
	csvService domain.CSVService
	logger     *logrus.Logger
}

// NewCSVHandler creates a new CSV handler
func NewCSVHandler(e *echo.Group, svc domain.CSVService, logger *logrus.Logger) {
	handler := &CSVHandler{
		csvService: svc,
		logger:     logger,
	}

	e.POST("/upload", handler.UploadCSV)

	e.GET("/jobs", handler.GetUserJobs)
	e.GET("/jobs/:job_id", handler.GetJobDetails)
	e.GET("/jobs/:job_id/progress", handler.GetJobProgress)
	e.GET("/jobs/:job_id/stream", handler.StreamProgress)
}

// UploadCSV handles multiple CSV file uploads
// @Summary Upload and process CSV files
// @Description Upload multiple CSV files for concurrent processing
// @Tags CSV
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "CSV files to upload"
// @Success 200 {object} domain.CSVUploadResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security ApiKeyAuth
// @Router /csv/upload [post]
func (h *CSVHandler) UploadCSV(c echo.Context) error {
	// Get user ID from context (set by JWT middleware)

	// Parse multipart form with memory limit (32MB)
	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		h.logger.WithError(err).Error("Failed to parse multipart form")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Failed to parse form data",
		})
	}

	// Get files from form
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get multipart form")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Failed to get form files",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "No files provided",
		})
	}

	// Validate file count (max 10 files)
	if len(files) > 10 {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Maximum 10 files allowed",
		})
	}

	// Validate file types
	for _, file := range files {
		if !isCSVFile(file.Filename) {
			return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Message: "Only CSV files are allowed",
			})
		}
	}
	// Process CSV files

	response, err := h.csvService.UploadAndProcessCSV(c.Request().Context(), files)
	if err != nil {
		h.logger.WithError(err).Error("Failed to processs CSV files")

		switch err {
		case domain.ErrBadParamInput:
			return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Message: "Failed to process CSV files",
			})
		}
	}

	return c.JSON(http.StatusOK, response)

}

// GetJobProgress retrieves the progress of a specific CSV processing job
// @Summary Get CSV job progress
// @Description Get real-time progress of a CSV processing job
// @Tags CSV
// @Produce json
// @Param job_id path string true "CSV Job ID"
// @Success 200 {object} domain.CSVProcessingProgress
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security ApiKeyAuth
// @Router /csv/jobs/{job_id}/progress [get]
func (h *CSVHandler) GetJobProgress(c echo.Context) error {
	jobID := c.Param("job_id")
	if jobID == "" {
		h.logger.Error("Job ID is required")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Job ID is required",
		})
	}

	jid, err := uuid.Parse(jobID)
	if err != nil {
		h.logger.WithError(err).WithField("job_id", jobID).Error("Invalid job ID format")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid job ID",
		})
	}

	progress, err := h.csvService.GetJobProgress(c.Request().Context(), jid)
	if err != nil {
		h.logger.WithError(err).WithField("job_id", jobID).Error("Failed to get job progress")

		switch err {
		case domain.ErrCSVJobNotFound:
			return c.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: "CSV job not found",
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Message: "Failed to get job progress",
			})
		}
	}

	return c.JSON(http.StatusOK, progress)
}

// GetUserJobs retrieves all CSV jobs for the authenticated user
// @Summary Get user CSV jobs
// @Description Get all CSV processing jobs for the authenticated user
// @Tags CSV
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.CSVJob}
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security ApiKeyAuth
// @Router /csv/jobs [get]
func (h *CSVHandler) GetUserJobs(c echo.Context) error {
	// Get user ID from context

	// Parse pagination parameters
	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	jobs, err := h.csvService.GetUserJobs(c.Request().Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user jobs")
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to get user jobs",
		})
	}

	// Apply pagination
	total := len(jobs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedJobs := jobs[start:end]

	response := domain.PaginatedResponse{
		Data: paginatedJobs,
		Pagination: domain.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: (total + limit - 1) / limit,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// GetJobDetails retrieves detailed information about a specific CSV job
// @Summary Get CSV job details
// @Description Get detailed information about a specific CSV processing job
// @Tags CSV
// @Produce json
func (h *CSVHandler) GetJobDetails(c echo.Context) error {
	jobID := c.Param("job_id")
	if jobID == "" {
		h.logger.Error("Job ID is required")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Job ID is required",
		})
	}

	jid, err := uuid.Parse(jobID)
	if err != nil {
		h.logger.WithError(err).WithField("job_id", jobID).Error("Invalid job ID format")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid job ID",
		})
	}

	progress, err := h.csvService.GetJobProgress(c.Request().Context(), jid)
	if err != nil {
		h.logger.WithError(err).WithField("job_id", jobID).Error("Failed to get job details")

		switch err {
		case domain.ErrCSVJobNotFound:
			return c.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: "CSV job not found",
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Message: "Failed to get job details",
			})
		}
	}

	return c.JSON(http.StatusOK, progress)
}
func (h *CSVHandler) StreamProgress(c echo.Context) error {
	jobID := c.Param("job_id")
	if jobID == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Job ID is required",
		})
	}

	jid, err := uuid.Parse(jobID)
	if err != nil {
		h.logger.WithError(err).WithField("job_id", jobID).Error("Invalid job ID format")
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid job ID",
		})
	}

	// Check if job exists
	_, err = h.csvService.GetJobProgress(c.Request().Context(), jid)
	if err != nil {
		if err == domain.ErrCSVJobNotFound {
			return c.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: "CSV job not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to get job details",
		})
	}

	// Set headers for SSE
	res := c.Response()
	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(http.StatusOK)

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		h.logger.Error("Streaming not supported")
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Streaming not supported",
		})
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			progress, err := h.csvService.GetJobProgress(c.Request().Context(), jid)
			if err != nil {
				h.logger.WithError(err).Error("Failed to get progress during streaming")
				return nil
			}

			// Marshal progress to JSON and write as SSE data
			data, err := json.Marshal(progress)
			if err != nil {
				h.logger.WithError(err).Error("Failed to marshal progress")
				return nil
			}

			if _, err := fmt.Fprintf(res.Writer, "event: progress\n"); err != nil {
				return nil
			}
			if _, err := fmt.Fprintf(res.Writer, "data: %s\n\n", data); err != nil {
				return nil
			}
			flusher.Flush()

			// Stop streaming if job is completed or failed
			if progress.Status == domain.CSVJobStatusCompleted || progress.Status == domain.CSVJobStatusFailed {
				return nil
			}
		}
	}
}

// isCSVFile checks if the filename has a CSV extension
func isCSVFile(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".csv"
}
