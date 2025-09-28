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

func TestUserCRUD_E2E(t *testing.T) {
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

	// Now start the test server
	kit.Start(t)

	//user := creUser.Data
	createReqUser := domain.CreateUserRequest{
		Name:     "John Doe Alex",
		Email:    "johnalex11@example.com",
		Password: "Password1234",
	}

	type CreateTypeUser domain.ResponseSingleData[domain.User]
	creUser, _ := doRequest[CreateTypeUser](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/users",
		createReqUser,
	)
	//end create user first
	userCreated := creUser.Data
	// Login to get JWT token using seeded user
	loginReq := domain.LoginRequest{
		Email:    userCreated.Email,
		Password: "Password1234",
	}
	type LoginType domain.ResponseSingleData[domain.LoginResponse]
	loginResp, code := doRequest[LoginType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/auth/login",
		loginReq,
	)
	require.Equal(t, http.StatusOK, code)
	token := loginResp.Data.Token

	// Get a seeded user for testing operations (since user creation requires admin privileges)
	type GetUsersType domain.ResponseMultipleData[domain.User]
	usersResp, code := doRequest[GetUsersType](
		t, http.MethodGet,
		kit.BaseURL+"/api/v1/users",
		nil,
		map[string]string{"Authorization": "Bearer " + token},
	)
	require.Equal(t, http.StatusOK, code)
	require.NotEmpty(t, usersResp.Data, "Should have seeded users")

	// Use the first seeded user for testing
	user := usersResp.Data[0]

	// Get
	type GetType domain.ResponseSingleData[domain.User]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
		map[string]string{"Authorization": "Bearer " + token},
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, user.ID, getE.Data.ID)

	// Update
	updPayload := domain.User{
		Name:  "Jane Doe",
		Email: "jane@example45.com",
	}
	type UpdType domain.ResponseSingleData[domain.User]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		updPayload,
		map[string]string{"Authorization": "Bearer " + token},
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "Jane Doe", updE.Data.Name)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Get after delete
	type ErrType domain.ResponseSingleData[domain.Empty]
	errE, code := doRequest[ErrType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
		map[string]string{"Authorization": "Bearer " + token},
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "User not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from users where id = $1", user.ID)
	require.NoError(t, err)
}
