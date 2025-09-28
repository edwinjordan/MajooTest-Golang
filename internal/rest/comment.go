package rest

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/internal/logging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *domain.CreateCommentRequest) (*domain.Comment, error)
	GetCommentList(ctx context.Context, filter *domain.CommentFilter) ([]domain.Comment, error)
	GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	UpdateComment(ctx context.Context, id uuid.UUID, comment *domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

type CommentHandler struct {
	Service CommentService
}

func NewCommentHandler(e *echo.Group, svc CommentService) {
	handler := &CommentHandler{
		Service: svc,
	}
	//commentsGroup := e.Group("/comments")
	e.GET("", handler.GetCommentList)
	e.GET("/:id", handler.GetComment)
	e.POST("", handler.CreateComment)
	e.PUT("/:id", handler.UpdateComment)
	e.DELETE("/:id", handler.DeleteComment)
}

// GetComments godoc
// @Summary List comments
// @Description Get all comments
// @Tags comments
// @Produce  json
// @Success 200 {array} domain.Comment
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /comments [get]
func (h *CommentHandler) GetCommentList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.CommentFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind comments filter", slog.String("error", err.Error()))
	}

	comments, err := h.Service.GetCommentList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_comments_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list comments: " + err.Error(),
		})
	}
	if comments == nil {
		comments = []domain.Comment{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Comment]{
		Data:    comments,
		Code:    http.StatusOK,
		Message: "Successfully retrieve comments list",
	})
}

// GetComment godoc
// @Summary Get post by ID
// @Description Get post details by ID
// @Tags comments
// @Produce  json
// @Param        id   path      int  true  "Post ID"
// @Success 200 {array} domain.Comment
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(c echo.Context) error {
	tracer := otel.Tracer("http.handler.comments")
	ctx, span := tracer.Start(c.Request().Context(), "GetCommentHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid comment ID format",
		})
	}

	span.SetAttributes(attribute.String("comment.id", id.String()))
	comment, err := h.Service.GetComment(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Message: "Comment not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		logging.LogError(ctx, err, "get_comment")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get comment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Comment]{
		Data:    *comment,
		Code:    http.StatusOK,
		Message: "Successfully retrieved comment",
	})
}

// CreateComments godoc
// @Summary Create comment
// @Description create a new comment
// @Tags comments
// @Accept  json
// @Produce  json
// @Param   comment  body  domain.CreateCommentRequest  true  "Comment data"
// @Success 201 {object} domain.CreateCommentRequest
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Router /comments [post]
func (h *CommentHandler) CreateComment(c echo.Context) error {
	var comment domain.CreateCommentRequest
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdComment, err := h.Service.CreateComment(ctx, &comment)
	if err != nil {
		logging.LogError(ctx, err, "create_comment")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create comment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Comment]{
		Data:    *createdComment,
		Code:    http.StatusCreated,
		Message: "Comment successfully created",
	})
}

// UpdateComment godoc
// @Summary Update comment
// @Description update an existing comment entry by ID
// @Tags comments
// @Accept  json
// @Produce  json
// @Param   id    path  string             true  "Comment ID"
// @Param   post  body  domain.UpdateCommentRequest  true  "Updated comment data"
// @Success 200 {object} domain.Comment
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid post ID format",
		})
	}

	var comment domain.Comment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedComment, err := h.Service.UpdateComment(ctx, id, &comment)
	if err != nil {
		logging.LogError(ctx, err, "update_comment")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update comment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Comment]{
		Data:    *updatedComment,
		Code:    http.StatusOK,
		Message: "Comment successfully updated",
	})
}

// DeleteComment godoc
// @Summary Delete comment
// @Description delete an existing comment entry by ID
// @Tags comments
// @Produce  json
// @Param   id   path  string  true  "Comment ID"
// @Success 204 {object} nil
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid comment ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteComment(ctx, id); err != nil {
		logging.LogError(ctx, err, "delete_comment")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete comment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Message: "User successfully deleted",
	})
}
