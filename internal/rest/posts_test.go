package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/internal/repository/postgres"
	"github.com/edwinjordan/MajooTest-Golang/internal/rest"
	"github.com/edwinjordan/MajooTest-Golang/internal/rest/middleware"
	"github.com/edwinjordan/MajooTest-Golang/service"
	"github.com/stretchr/testify/require"
)

func TestPostsCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Create API v1 group
	apiV1 := kit.Echo.Group("/api/v1")

	// Wire auth routes (no authentication required)
	authRepo := postgres.NewAuthRepository(kit.DB)
	authSvc := service.NewAuthService(authRepo)
	rest.NewAuthHandler(apiV1.Group("/auth"), authSvc)

	// Wire user routes (no authentication for user creation, but protected for other operations)
	userRepo := postgres.NewUserRepository(kit.DB)
	userSvc := service.NewUserService(userRepo)
	rest.NewUserHandler(apiV1.Group("/users"), userSvc)

	// Wire posts routes (with authentication)
	postsRepo := postgres.NewPostsRepository(kit.DB)
	postsSvc := service.NewPostsService(postsRepo)
	postsGroup := apiV1.Group("/posts", middleware.ValidateUserToken())
	rest.NewPostsHandler(postsGroup, postsSvc)

	// Now start the test server
	kit.Start(t)

	//user := creUser.Data
	createReqUser := domain.CreateUserRequest{
		Name:     "John Doe Alex",
		Email:    "johnalex@example.com",
		Password: "Password1234",
	}

	type CreateTypeUser domain.ResponseSingleData[domain.User]
	creUser, _ := doRequest[CreateTypeUser](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/users",
		createReqUser,
	)
	//end create user first
	user := creUser.Data
	// Login to obtain JWT token first
	type LoginType domain.ResponseSingleData[domain.LoginResponse]
	loginReq := domain.LoginRequest{
		Email:    user.Email,     // Use the created user's email
		Password: "Password1234", // Use the seeded user's password
	}
	loginRes, code := doRequest[LoginType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/auth/login",
		loginReq,
	)
	require.Equal(t, http.StatusOK, code)
	token := "Bearer " + loginRes.Data.Token
	// end login
	// Create auth headers map
	authHeaders := map[string]string{
		"Authorization": token,
	}
	// Create
	createReq := domain.CreatePostsRequest{
		Title:   "Test Content",
		Content: "Ini adalah content",
	}
	type CreateTypePost domain.ResponseSingleData[domain.Posts]
	cre, code := doRequest[CreateTypePost](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/posts",
		createReq,
		authHeaders,
	)
	require.Equal(t, http.StatusCreated, code)
	post := cre.Data
	require.NotEmpty(t, post.ID)

	// Get
	type GetType domain.ResponseSingleData[domain.Posts]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/posts/%s", kit.BaseURL, post.ID),
		nil,
		authHeaders,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, post.ID, getE.Data.ID)

	// Update
	updPayload := domain.Posts{
		Title:   "Jane Doe",
		Content: "Ini adalah content",
	}
	type UpdType domain.ResponseSingleData[domain.Posts]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/posts/%s", kit.BaseURL, post.ID),
		updPayload,
		authHeaders,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "Jane Doe", updE.Data.Title)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/posts/%s", kit.BaseURL, post.ID),
		nil,
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Get after delete
	type ErrType domain.ResponseSingleData[domain.Empty]
	errE, code := doRequest[ErrType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/posts/%s", kit.BaseURL, post.ID),
		nil,
		map[string]string{"Authorization": "Bearer " + token},
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "Post not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from posts where id = $1", post.ID)
	require.NoError(t, err)
}
