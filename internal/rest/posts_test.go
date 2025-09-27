package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/internal/repository/postgres"
	"github.com/edwinjordan/MajooTest-Golang/internal/rest"
	"github.com/edwinjordan/MajooTest-Golang/service"
	"github.com/stretchr/testify/require"
)

func TestPostsCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	postsRepo := postgres.NewPostsRepository(kit.DB)
	postsSvc := service.NewPostsService(postsRepo)
	rest.NewPostsHandler(kit.Echo.Group("/api/v1"), postsSvc)

	// Now start the test server
	kit.Start(t)

	// Create
	createReq := domain.CreatePostsRequest{
		Title:   "Test Content",
		Content: "Ini adalah content",
	}
	type CreateType domain.ResponseSingleData[domain.Posts]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/posts",
		createReq,
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
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "Post not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from posts where id = $1", post.ID)
	require.NoError(t, err)
}
