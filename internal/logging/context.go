package logging

import (
	"context"
	"log/slog"

	"github.com/edwinjordan/MajooTest-Golang/internal/rest/middleware"
)

// TopicInfo holds Topic context information for logging
type TopicInfo struct {
	ID        string `json:"Topic_id,omitempty"`
	Topicname string `json:"Topicname,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
}

const (
	TopicContextKey = "Topic_info"
)

// WithTopicInfo adds Topic information to the context
func WithTopicInfo(ctx context.Context, TopicInfo *TopicInfo) context.Context {
	return context.WithValue(ctx, TopicContextKey, TopicInfo)
}

// GetTopicInfo extracts Topic information from context
func GetTopicInfo(ctx context.Context) *TopicInfo {
	if TopicInfo, ok := ctx.Value(TopicContextKey).(*TopicInfo); ok {
		return TopicInfo
	}
	return nil
}

// NewContextualLogger creates a logger with request ID and Topic context
func NewContextualLogger(ctx context.Context) *slog.Logger {
	logger := slog.Default()

	// Add request ID if available
	requestID := middleware.GetRequestID(ctx)
	if requestID != "" {
		logger = logger.With(slog.String("request_id", requestID))
	}

	// Add Topic information if available
	TopicInfo := GetTopicInfo(ctx)
	if TopicInfo != nil {
		attrs := make([]any, 0, 4)
		if TopicInfo.ID != "" {
			attrs = append(attrs, slog.String("Topic_id", TopicInfo.ID))
		}
		if TopicInfo.Topicname != "" {
			attrs = append(attrs, slog.String("Topicname", TopicInfo.Topicname))
		}
		if TopicInfo.Email != "" {
			attrs = append(attrs, slog.String("Topic_email", TopicInfo.Email))
		}
		if TopicInfo.Role != "" {
			attrs = append(attrs, slog.String("Topic_role", TopicInfo.Role))
		}
		if len(attrs) > 0 {
			logger = logger.With(attrs...)
		}
	}

	return logger
}

// LoggerWithFields creates a logger with additional fields
func LoggerWithFields(ctx context.Context, fields ...any) *slog.Logger {
	logger := NewContextualLogger(ctx)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return logger
}

// Security logging functions
func LogSecurityEvent(ctx context.Context, event string, details ...any) {
	logger := NewContextualLogger(ctx)
	args := []any{slog.String("security_event", event)}
	args = append(args, details...)
	logger.Warn("Security Event", args...)
}

func LogAuthAttempt(ctx context.Context, Topicname string, success bool, reason string) {
	logger := NewContextualLogger(ctx)
	logger.Info("Authentication Attempt",
		slog.String("Topicname", Topicname),
		slog.Bool("success", success),
		slog.String("reason", reason),
	)
}

func LogDataAccess(ctx context.Context, resource string, action string, result string) {
	logger := NewContextualLogger(ctx)
	logger.Info("Data Access",
		slog.String("resource", resource),
		slog.String("action", action),
		slog.String("result", result),
	)
}

// Performance logging functions
func LogPerformance(ctx context.Context, operation string, duration int64, metadata ...any) {
	logger := NewContextualLogger(ctx)
	args := []any{
		slog.String("operation", operation),
		slog.Int64("duration_ms", duration),
	}
	args = append(args, metadata...)
	logger.Info("Performance Metric", args...)
}

// Business logic logging functions
func LogBusinessEvent(ctx context.Context, event string, entityType string, entityID string, details ...any) {
	logger := NewContextualLogger(ctx)
	args := []any{
		slog.String("business_event", event),
		slog.String("entity_type", entityType),
		slog.String("entity_id", entityID),
	}
	args = append(args, details...)
	logger.Info("Business Event", args...)
}

// General info logging with context
func LogInfo(ctx context.Context, message string, details ...any) {
	logger := NewContextualLogger(ctx)
	logger.Info(message, details...)
}

// General warn logging with context
func LogWarn(ctx context.Context, message string, details ...any) {
	logger := NewContextualLogger(ctx)
	logger.Info(message, details...)
}

// Error message logging with context (when you have a message but no error object)
func LogErrorMessage(ctx context.Context, message string, details ...any) {
	logger := NewContextualLogger(ctx)
	logger.Error(message, details...)
}

// Error logging with context
func LogError(ctx context.Context, err error, operation string, details ...any) {
	logger := NewContextualLogger(ctx)
	args := []any{
		slog.String("error", err.Error()),
		slog.String("operation", operation),
	}
	args = append(args, details...)
	logger.Error("Operation Failed", args...)
}

func LogErrorWithStackTrace(ctx context.Context, err error, operation string, stackTrace string, details ...any) {
	logger := NewContextualLogger(ctx)
	args := []any{
		slog.String("error", err.Error()),
		slog.String("operation", operation),
		slog.String("stack_trace", stackTrace),
	}
	args = append(args, details...)
	logger.Error("Operation Failed with Stack Trace", args...)
}
