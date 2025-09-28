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

func TestCommentCRUD_E2E(t *testing.T) {
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

	// Wire comment routes (with authentication)
	commentRepo := postgres.NewCommentRepository(kit.DB)
	commentSvc := service.NewCommentService(commentRepo)
	commentGroup := apiV1.Group("/comments", middleware.ValidateUserToken())
	rest.NewCommentHandler(commentGroup, commentSvc)

	// Now start the test server
	kit.Start(t)
	//create user first
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

	user := creUser.Data
	//end create user first

	// Login to obtain JWT token first
	type LoginType domain.ResponseSingleData[domain.LoginResponse]
	loginReq := domain.LoginRequest{
		Email:    "johnalex@example.com", // Use the created user's email
		Password: "Password1234",         // Use the created user's password
	}
	loginRes, code := doRequest[LoginType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/auth/login",
		loginReq,
	)
	require.Equal(t, http.StatusOK, code)
	token := "Bearer " + loginRes.Data.Token

	// Create auth headers map
	authHeaders := map[string]string{
		"Authorization": token,
	}

	//create post first (with authentication)
	createReqPost := domain.CreatePostsRequest{
		Title:   "Double Post",
		Content: "This is the content of the first post",
	}
	type CreateTypePost domain.ResponseSingleData[domain.Posts]
	crePost, _ := doRequest[CreateTypePost](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/posts",
		createReqPost,
		authHeaders,
	)
	posts := crePost.Data
	//end create post first

	// Create comment with authentication

	createReq := domain.CreateCommentRequest{
		UserID: user.ID,
		PostID: posts.ID,
		Body:   "This is a comment",
	}
	type CreateType domain.ResponseSingleData[domain.Comment]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/comments",
		createReq,
		authHeaders,
	)
	require.Equal(t, http.StatusCreated, code)
	comment := cre.Data
	require.NotEmpty(t, comment.ID)

	// Get
	type GetType domain.ResponseSingleData[domain.Comment]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/comments/%s", kit.BaseURL, comment.ID),
		nil,
		authHeaders,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, comment.ID, getE.Data.ID)

	// Update

	//create user first
	createReqUserUpdate := domain.CreateUserRequest{
		Name:     "John Doe Asa",
		Email:    "john@exampleasa.com",
		Password: "Password1234",
	}
	type CreateTypeUserUpdate domain.ResponseSingleData[domain.User]
	creUserUpdate, _ := doRequest[CreateTypeUserUpdate](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/users",
		createReqUserUpdate,
	)
	userUpdate := creUserUpdate.Data
	//end create user first

	//create post first (with authentication)
	createReqPostUpdate := domain.CreatePostsRequest{
		Title:   "Second Post",
		Content: "This is the content of the second post",
	}
	type CreateTypePostUpdate domain.ResponseSingleData[domain.Posts]
	crePostUpdate, _ := doRequest[CreateTypePostUpdate](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/posts",
		createReqPostUpdate,
		authHeaders,
	)
	postsUpdate := crePostUpdate.Data
	//end create post first

	updPayload := domain.Comment{
		UserID: userUpdate.ID,
		PostID: postsUpdate.ID,
		Body:   "This is an updated comment",
	}
	type UpdType domain.ResponseSingleData[domain.Comment]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/comments/%s", kit.BaseURL, comment.ID),
		updPayload,
		authHeaders,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "This is an updated comment", updE.Data.Body)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/comments/%s", kit.BaseURL, comment.ID),
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
		fmt.Sprintf("%s/api/v1/comments/%s", kit.BaseURL, comment.ID),
		nil,
		authHeaders,
	)

	//get after delete user
	_, err = kit.DB.Exec(context.Background(), "DELETE from users where id = $1", user.ID)
	require.NoError(t, err)

	//get after delete post
	_, err = kit.DB.Exec(context.Background(), "DELETE from posts where id = $1", posts.ID)
	require.NoError(t, err)

	//get after delete user update
	_, err = kit.DB.Exec(context.Background(), "DELETE from users where id = $1", userUpdate.ID)
	require.NoError(t, err)

	//get after delete post update
	_, err = kit.DB.Exec(context.Background(), "DELETE from posts where id = $1", postsUpdate.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "Comment not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from comments where id = $1", comment.ID)
	require.NoError(t, err)
}
