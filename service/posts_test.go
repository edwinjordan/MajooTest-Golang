package service_test

import (
	"context"
	"errors"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/service"
	"github.com/edwinjordan/MajooTest-Golang/service/mocks"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostsService_CreatePosts(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)

	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	req := &domain.CreatePostsRequest{
		Title:   "Test Post",
		Content: "This is a test post content.",
		Slug:    "test-post",
	}
	expectedPosts := &domain.Posts{
		ID:      uuid.New().String(),
		Title:   "Test Post",
		Content: "This is a test post content.",
		Slug:    "test-post",
	}

	t.Run("Successfully creates a post", func(t *testing.T) {
		mockPostsRepo.On("CreatePosts", mock.Anything, req).Return(expectedPosts, nil).Once()

		posts, err := postsService.CreatePosts(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Equal(t, expectedPosts.ID, posts.ID)
		assert.Equal(t, expectedPosts.Title, posts.Title)
		assert.Equal(t, expectedPosts.Content, posts.Content)
		assert.Equal(t, expectedPosts.Slug, posts.Slug)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("database error")
		mockPostsRepo.On("CreatePosts", mock.Anything, req).Return(nil, repoErr).Once()

		posts, err := postsService.CreatePosts(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, repoErr, err)

		mockPostsRepo.AssertExpectations(t)
	})
}

func TestPostsService_GetPosts(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)
	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	postsID := uuid.New()
	expectedPosts := &domain.Posts{
		ID:      postsID.String(),
		Title:   "Fetched Post",
		Content: "This is fetched post content.",
		Slug:    "fetched-post",
	}

	t.Run("Successfully fetches a post", func(t *testing.T) {
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(expectedPosts, nil).Once()

		posts, err := postsService.GetPosts(ctx, postsID)

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Equal(t, expectedPosts.ID, posts.ID)
		assert.Equal(t, expectedPosts.Title, posts.Title)
		assert.Equal(t, expectedPosts.Content, posts.Content)
		assert.Equal(t, expectedPosts.Slug, posts.Slug)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("network error")
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, repoErr).Once()

		posts, err := postsService.GetPosts(ctx, postsID)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, repoErr, err)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when post not found in repository", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, nil).Once()

		posts, err := postsService.GetPosts(ctx, postsID)

		assert.NoError(t, err)
		assert.Nil(t, posts)

		mockPostsRepo.AssertExpectations(t)
	})
}

func TestPostsService_UpdatePosts(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)
	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	postsID := uuid.New()
	existingPosts := &domain.Posts{
		ID:      postsID.String(),
		Title:   "Old Title",
		Content: "old content",
		Slug:    "old-title",
	}
	updateReq := &domain.Posts{
		Title:   "New Title",
		Content: "new content",
		Slug:    "new-title",
	}

	t.Run("Successfully updates a post", func(t *testing.T) {
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()

		expectedUpdatedPosts := &domain.Posts{
			ID:      postsID.String(),
			Title:   updateReq.Title,
			Content: updateReq.Content,
			Slug:    updateReq.Slug,
		}
		mockPostsRepo.On("UpdatePosts", mock.Anything, postsID, expectedUpdatedPosts).Return(expectedUpdatedPosts, nil).Once()

		posts, err := postsService.UpdatePosts(ctx, postsID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Equal(t, expectedUpdatedPosts.Title, posts.Title)
		assert.Equal(t, expectedUpdatedPosts.Content, posts.Content)
		assert.Equal(t, expectedUpdatedPosts.Slug, posts.Slug)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrPostsNotFound if posts does not exist", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, nil).Once()

		posts, err := postsService.UpdatePosts(ctx, postsID, updateReq)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, posts)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetPosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("get posts repo error")
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, repoErr).Once()

		posts, err := postsService.UpdatePosts(ctx, postsID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, repoErr, err)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if UpdatePosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()

		repoErr := errors.New("update posts repo error")
		expectedUpdatedPosts := &domain.Posts{
			ID:      postsID.String(),
			Title:   updateReq.Title,
			Content: updateReq.Content,
			Slug:    updateReq.Slug,
		}
		mockPostsRepo.On("UpdatePosts", mock.Anything, postsID, expectedUpdatedPosts).Return(nil, repoErr).Once()

		posts, err := postsService.UpdatePosts(ctx, postsID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, repoErr, err)

		mockPostsRepo.AssertExpectations(t)
	})
}

func TestPostService_DeletePost(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)
	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	postsID := uuid.New()
	existingPosts := &domain.Posts{
		ID:      postsID.String(),
		Title:   "User to delete",
		Content: "delete@example.com",
		Slug:    "user-to-delete",
	}

	t.Run("Successfully deletes a post", func(t *testing.T) {
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()
		mockPostsRepo.On("DeletePosts", mock.Anything, postsID).Return(nil).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.NoError(t, err)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrPostsNotFound if posts does not exist", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, nil).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetPosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("get posts repo error during delete")
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, repoErr).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeletePosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()
		repoErr := errors.New("delete posts repo error")
		mockPostsRepo.On("DeletePosts", mock.Anything, postsID).Return(repoErr).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockPostsRepo.AssertExpectations(t)
	})
}

func TestPostsService_DeletePost(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)
	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	postsID := uuid.New()
	existingPosts := &domain.Posts{
		ID:      postsID.String(),
		Title:   "User to delete",
		Content: "delete@example.com",
		Slug:    "user-to-delete",
	}

	t.Run("Successfully deletes a posts", func(t *testing.T) {
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()
		mockPostsRepo.On("DeletePosts", mock.Anything, postsID).Return(nil).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.NoError(t, err)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrPostsNotFound if posts does not exist", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, nil).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetPosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("get posts repo error during delete")
		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(nil, repoErr).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeletePosts fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPosts", mock.Anything, postsID).Return(existingPosts, nil).Once()
		repoErr := errors.New("delete posts repo error")
		mockPostsRepo.On("DeletePosts", mock.Anything, postsID).Return(repoErr).Once()

		err := postsService.DeletePosts(ctx, postsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockPostsRepo.AssertExpectations(t)
	})
}

func TestPostsService_GetPostsList(t *testing.T) {
	mockPostsRepo := new(mocks.PostsRepository)
	postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	filter := &domain.PostsFilter{
		Search: "test",
	}
	expectedPosts := []domain.Posts{
		{ID: uuid.New().String(), Title: "Test Post One", Content: "Content One", Slug: "test-post-one"},
		{ID: uuid.New().String(), Title: "Another Test Post", Content: "Content Two", Slug: "another-test-post"},
	}

	t.Run("Successfully fetches post list", func(t *testing.T) {
		mockPostsRepo.On("GetPostsList", mock.Anything, filter).Return(expectedPosts, nil).Once()

		posts, err := postsService.GetPostsList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Len(t, posts, 2)
		assert.Equal(t, expectedPosts[0].Title, posts[0].Title)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no posts found", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		mockPostsRepo.On("GetPostsList", mock.Anything, filter).Return([]domain.Posts{}, nil).Once()

		posts, err := postsService.GetPostsList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Len(t, posts, 0)

		mockPostsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockPostsRepo = new(mocks.PostsRepository)
		postsService = service.NewPostsService(mockPostsRepo)

		repoErr := errors.New("get posts list database error")
		mockPostsRepo.On("GetPostsList", mock.Anything, filter).Return(nil, repoErr).Once()

		posts, err := postsService.GetPostsList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, repoErr, err)

		mockPostsRepo.AssertExpectations(t)
	})
}
