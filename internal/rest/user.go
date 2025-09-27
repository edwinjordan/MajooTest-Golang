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

type UserService interface {
	CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(e *echo.Group, svc UserService) {
	handler := &UserHandler{
		Service: svc,
	}
	usersGroup := e.Group("/users")
	usersGroup.GET("", handler.GetUserList)
	usersGroup.GET("/:id", handler.GetUser)
	usersGroup.POST("", handler.CreateUser)
	usersGroup.PUT("/:id", handler.UpdateUser)
	usersGroup.DELETE("/:id", handler.DeleteUser)
}

// GetUser godoc
// @Summary List user
// @Description Get all user
// @Tags user
// @Produce  json
// @Success 200 {array} domain.User
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) GetUserList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.UserFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind user filter", slog.String("error", err.Error()))
	}

	users, err := h.Service.GetUserList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_user_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list users: " + err.Error(),
		})
	}
	if users == nil {
		users = []domain.User{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.User]{
		Data:    users,
		Code:    http.StatusOK,
		Message: "Successfully retrieve user list",
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags user
// @Produce  json
// @Param        id   path      int  true  "Account ID"
// @Success 200 {array} domain.User
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "GetUserHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	span.SetAttributes(attribute.String("user.id", id.String()))
	user, err := h.Service.GetUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Message: "User not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		logging.LogError(ctx, err, "get_user")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *user,
		Code:    http.StatusOK,
		Message: "Successfully retrieved user",
	})
}

// CreateUser godoc
// @Summary Create user
// @Description create a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param   user  body  domain.CreateUserRequest  true  "User data"
// @Success 201 {object} domain.CreateUserRequest
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var user domain.CreateUserRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdUser, err := h.Service.CreateUser(ctx, &user)
	if err != nil {
		logging.LogError(ctx, err, "create_user")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.User]{
		Data:    *createdUser,
		Code:    http.StatusCreated,
		Message: "User successfully created",
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description update an existing user entry by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param   id    path  string             true  "User ID"
// @Param   user  body  domain.UpdateUserRequest  true  "Updated user data"
// @Success 200 {object} domain.User
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedUser, err := h.Service.UpdateUser(ctx, id, &user)
	if err != nil {
		logging.LogError(ctx, err, "update_user")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *updatedUser,
		Code:    http.StatusOK,
		Message: "User successfully updated",
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
