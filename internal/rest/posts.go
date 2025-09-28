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

type PostsService interface {
	CreatePosts(ctx context.Context, post *domain.CreatePostsRequest) (*domain.Posts, error)
	GetPostsList(ctx context.Context, filter *domain.PostsFilter) ([]domain.Posts, error)
	GetPosts(ctx context.Context, id uuid.UUID) (*domain.Posts, error)
	UpdatePosts(ctx context.Context, id uuid.UUID, post *domain.Posts) (*domain.Posts, error)
	DeletePosts(ctx context.Context, id uuid.UUID) error
}

type PostsHandler struct {
	Service PostsService
}

func NewPostsHandler(e *echo.Group, svc PostsService) {
	handler := &PostsHandler{
		Service: svc,
	}
	//postsGroup := e.Group("/posts")
	e.GET("", handler.GetPostsList)
	e.GET("/:id", handler.GetPosts)
	e.POST("", handler.CreatePosts)
	e.PUT("/:id", handler.UpdatePosts)
	e.DELETE("/:id", handler.DeletePosts)
}

// GetPosts godoc
// @Summary List posts
// @Description Get all posts
// @Tags posts
// @Produce  json
// @Success 200 {array} domain.Posts
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /posts [get]
func (h *PostsHandler) GetPostsList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.PostsFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind posts filter", slog.String("error", err.Error()))
	}

	posts, err := h.Service.GetPostsList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_posts_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list posts: " + err.Error(),
		})
	}
	if posts == nil {
		posts = []domain.Posts{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Posts]{
		Data:    posts,
		Code:    http.StatusOK,
		Message: "Successfully retrieve posts list",
	})
}

// GetPosts godoc
// @Summary Get post by ID
// @Description Get post details by ID
// @Tags posts
// @Produce  json
// @Param        id   path      int  true  "Post ID"
// @Success 200 {array} domain.Posts
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /posts/{id} [get]
func (h *PostsHandler) GetPosts(c echo.Context) error {
	tracer := otel.Tracer("http.handler.posts")
	ctx, span := tracer.Start(c.Request().Context(), "GetUserHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid post ID format",
		})
	}

	span.SetAttributes(attribute.String("post.id", id.String()))
	post, err := h.Service.GetPosts(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Message: "Post not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		logging.LogError(ctx, err, "get_post")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get post: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Posts]{
		Data:    *post,
		Code:    http.StatusOK,
		Message: "Successfully retrieved user",
	})
}

// CreatePosts godoc
// @Summary Create post
// @Description create a new post
// @Tags posts
// @Accept  json
// @Produce  json
// @Param   post  body  domain.CreatePostsRequestSwagger  true  "Post data"
// @Success 201 {object} domain.CreatePostsRequest
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Router /posts [post]
func (h *PostsHandler) CreatePosts(c echo.Context) error {
	var post domain.CreatePostsRequest
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdPost, err := h.Service.CreatePosts(ctx, &post)
	if err != nil {
		logging.LogError(ctx, err, "create_post")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create post: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Posts]{
		Data:    *createdPost,
		Code:    http.StatusCreated,
		Message: "Post successfully created",
	})
}

// UpdatePosts godoc
// @Summary Update post
// @Description update an existing post entry by ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param   id    path  string             true  "Post ID"
// @Param   post  body  domain.UpdatePostsRequestSwagger  true  "Updated post data"
// @Success 200 {object} domain.Posts
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /posts/{id} [put]
func (h *PostsHandler) UpdatePosts(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid post ID format",
		})
	}

	var post domain.Posts
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedPost, err := h.Service.UpdatePosts(ctx, id, &post)
	if err != nil {
		logging.LogError(ctx, err, "update_post")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update post: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Posts]{
		Data:    *updatedPost,
		Code:    http.StatusOK,
		Message: "Post successfully updated",
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description delete an existing user entry by ID
// @Tags user
// @Produce  json
// @Param   id   path  string  true  "User ID"
// @Success 204 {object} nil
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteUser(ctx, id); err != nil {
		logging.LogError(ctx, err, "delete_user")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Message: "User successfully deleted",
	})
}

// DeletePost godoc
// @Summary Delete post
// @Description delete an existing post entry by ID
// @Tags posts
// @Produce  json
// @Param   id   path  string  true  "Post ID"
// @Success 204 {object} nil
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
func (h *PostsHandler) DeletePosts(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeletePosts(ctx, id); err != nil {
		logging.LogError(ctx, err, "delete_post")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete post: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Message: "User successfully deleted",
	})
}
