package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/edwinjordan/MajooTest-Golang/service"
	"github.com/edwinjordan/MajooTest-Golang/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment(t *testing.T) {
	mockCommentsRepo := new(mocks.CommentRepository)
	// mockUserRepo := new(mocks.UserRepository)
	// mockPostsRepo := new(mocks.PostsRepository)

	//userService := service.NewUserService(mockUserRepo)
	commentsService := service.NewCommentService(mockCommentsRepo)
	//postsService := service.NewPostsService(mockPostsRepo)

	ctx := context.Background()
	req := &domain.CreatePostsRequest{
		Title:   "Test Post",
		Content: "This is a test post content.",
		Slug:    "test-post",
	}

	expectedPosts := &domain.Posts{
		ID:      uuid.New().String(),
		Title:   req.Title,
		Content: req.Content,
		Slug:    req.Slug,
	}

	//set user
	reqUser := &domain.CreateUserRequest{
		Name:  "Test User",
		Email: "testuser@example.com",
	}

	expectedUser := &domain.User{
		ID:    uuid.New().String(),
		Name:  reqUser.Name,
		Email: reqUser.Email,
	}
	//end set user

	// use the post ID (as if returned from posts service) for the comment request
	reqComment := &domain.CreateCommentRequest{
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "This is a test comment content.",
	}

	// if later tests expect a variable named 'req' for posts creation, keep using 'req'
	// and use 'reqComment' when creating comments
	expectedComments := &domain.Comment{
		ID:     uuid.New().String(),
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "This is a test comment content.",
	}

	t.Run("Successfully creates a comment", func(t *testing.T) {
		mockCommentsRepo.On("CreateComment", mock.Anything, reqComment).Return(expectedComments, nil).Once()

		comments, err := commentsService.CreateComment(ctx, reqComment)

		assert.NoError(t, err)
		assert.NotNil(t, comments)
		assert.Equal(t, expectedComments.ID, comments.ID)
		assert.Equal(t, expectedComments.PostID, comments.PostID)
		assert.Equal(t, expectedComments.UserID, comments.UserID)
		assert.Equal(t, expectedComments.Body, comments.Body)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		repoErr := errors.New("database error")
		mockCommentsRepo.On("CreateComment", mock.Anything, reqComment).Return(nil, repoErr).Once()

		comments, err := commentsService.CreateComment(ctx, reqComment)

		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Equal(t, repoErr, err)

		mockCommentsRepo.AssertExpectations(t)
	})
}

func TestCommentService_GetComment(t *testing.T) {
	mockCommentsRepo := new(mocks.CommentRepository)
	commentsService := service.NewCommentService(mockCommentsRepo)

	ctx := context.Background()

	req := &domain.CreatePostsRequest{
		Title:   "Test Post",
		Content: "This is a test post content.",
		Slug:    "test-post",
	}

	expectedPosts := &domain.Posts{
		ID:      uuid.New().String(),
		Title:   req.Title,
		Content: req.Content,
		Slug:    req.Slug,
	}

	//set user
	reqUser := &domain.CreateUserRequest{
		Name:  "Test User",
		Email: "testuser@example.com",
	}

	expectedUser := &domain.User{
		ID:    uuid.New().String(),
		Name:  reqUser.Name,
		Email: reqUser.Email,
	}
	//end set user

	commentID := uuid.New()
	expectedComment := &domain.Comment{
		ID:     commentID.String(),
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "Fetched Comment",
	}

	t.Run("Successfully fetches a comment", func(t *testing.T) {
		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(expectedComment, nil).Once()

		comment, err := commentsService.GetComment(ctx, commentID)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, expectedComment.ID, comment.ID)
		assert.Equal(t, expectedComment.PostID, comment.PostID)
		assert.Equal(t, expectedComment.UserID, comment.UserID)
		assert.Equal(t, expectedComment.Body, comment.Body)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		repoErr := errors.New("network error")
		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(nil, repoErr).Once()

		comment, err := commentsService.GetComment(ctx, commentID)

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, repoErr, err)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when comment not found in repository", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(nil, nil).Once()

		comment, err := commentsService.GetComment(ctx, commentID)

		assert.NoError(t, err)
		assert.Nil(t, comment)

		mockCommentsRepo.AssertExpectations(t)
	})
}

func TestCommentService_UpdateComment(t *testing.T) {
	mockCommentsRepo := new(mocks.CommentRepository)
	commentsService := service.NewCommentService(mockCommentsRepo)

	ctx := context.Background()

	req := &domain.CreatePostsRequest{
		Title:   "Test Post",
		Content: "This is a test post content.",
		Slug:    "test-post",
	}

	expectedPosts := &domain.Posts{
		ID:      uuid.New().String(),
		Title:   req.Title,
		Content: req.Content,
		Slug:    req.Slug,
	}

	//set user
	reqUser := &domain.CreateUserRequest{
		Name:  "Test User",
		Email: "testuser@example.com",
	}

	expectedUser := &domain.User{
		ID:    uuid.New().String(),
		Name:  reqUser.Name,
		Email: reqUser.Email,
	}
	//end set user
	commentID := uuid.New()
	existingComment := &domain.Comment{
		ID:     commentID.String(),
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "Old Comment",
	}
	updateReq := &domain.Comment{
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "Updated Comment",
	}

	t.Run("Successfully updates a comment", func(t *testing.T) {
		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(existingComment, nil).Once()

		expectedUpdatedComment := &domain.Comment{
			ID:     commentID.String(),
			PostID: expectedPosts.ID,
			UserID: expectedUser.ID,
			Body:   updateReq.Body,
		}
		mockCommentsRepo.On("UpdateComment", mock.Anything, commentID, expectedUpdatedComment).Return(expectedUpdatedComment, nil).Once()

		comment, err := commentsService.UpdateComment(ctx, commentID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, expectedUpdatedComment.Body, comment.Body)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrCommentNotFound if comment does not exist", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(nil, nil).Once()

		comment, err := commentsService.UpdateComment(ctx, commentID, updateReq)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, comment)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetComment fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		repoErr := errors.New("get comment repo error")
		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(nil, repoErr).Once()

		comment, err := commentsService.UpdateComment(ctx, commentID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, repoErr, err)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if UpdateComment fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetComment", mock.Anything, commentID).Return(existingComment, nil).Once()

		repoErr := errors.New("update comment repo error")
		expectedUpdatedComment := &domain.Comment{
			ID:     commentID.String(),
			PostID: updateReq.PostID,
			UserID: updateReq.UserID,
			Body:   updateReq.Body,
		}
		mockCommentsRepo.On("UpdateComment", mock.Anything, commentID, expectedUpdatedComment).Return(nil, repoErr).Once()

		comment, err := commentsService.UpdateComment(ctx, commentID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, repoErr, err)

		mockCommentsRepo.AssertExpectations(t)
	})
}

func TestCommentService_DeleteComment(t *testing.T) {
	mockCommentsRepo := new(mocks.CommentRepository)
	commentsService := service.NewCommentService(mockCommentsRepo)

	ctx := context.Background()

	expectedPosts := &domain.Posts{
		ID: uuid.New().String(),
	}

	expectedUser := &domain.User{
		ID: uuid.New().String(),
	}

	userID := uuid.New()
	existingUser := &domain.Comment{
		ID:     userID.String(),
		PostID: expectedPosts.ID,
		UserID: expectedUser.ID,
		Body:   "Some Comment",
	}

	t.Run("Successfully deletes a comment", func(t *testing.T) {
		mockCommentsRepo.On("GetComment", mock.Anything, userID).Return(existingUser, nil).Once()
		mockCommentsRepo.On("DeleteComment", mock.Anything, userID).Return(nil).Once()

		err := commentsService.DeleteComment(ctx, userID)

		assert.NoError(t, err)
		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrCommentNotFound if comment does not exist", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetComment", mock.Anything, userID).Return(nil, nil).Once()

		err := commentsService.DeleteComment(ctx, userID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeleteComment fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		repoErr := errors.New("get comment repo error during delete")
		mockCommentsRepo.On("GetComment", mock.Anything, userID).Return(nil, repoErr).Once()

		err := commentsService.DeleteComment(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeleteComment fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetComment", mock.Anything, userID).Return(existingUser, nil).Once()
		repoErr := errors.New("delete comment repo error")
		mockCommentsRepo.On("DeleteComment", mock.Anything, userID).Return(repoErr).Once()

		err := commentsService.DeleteComment(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockCommentsRepo.AssertExpectations(t)
	})
}

func TestCommentService_GetCommentList(t *testing.T) {
	mockCommentsRepo := new(mocks.CommentRepository)
	commentsService := service.NewCommentService(mockCommentsRepo)

	ctx := context.Background()
	filter := &domain.CommentFilter{
		Search: "test",
	}
	expectedComments := []domain.Comment{
		{ID: uuid.New().String(), PostID: uuid.New().String(), UserID: uuid.New().String(), Body: "Test Comment One"},
		{ID: uuid.New().String(), PostID: uuid.New().String(), UserID: uuid.New().String(), Body: "Another Test Comment"},
	}

	t.Run("Successfully fetches comment list", func(t *testing.T) {
		mockCommentsRepo.On("GetCommentList", mock.Anything, filter).Return(expectedComments, nil).Once()

		comments, err := commentsService.GetCommentList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, comments)
		assert.Len(t, comments, 2)
		assert.Equal(t, expectedComments[0].Body, comments[0].Body)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no comments found", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		mockCommentsRepo.On("GetCommentList", mock.Anything, filter).Return([]domain.Comment{}, nil).Once()

		comments, err := commentsService.GetCommentList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, comments)
		assert.Len(t, comments, 0)

		mockCommentsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockCommentsRepo = new(mocks.CommentRepository)
		commentsService = service.NewCommentService(mockCommentsRepo)

		repoErr := errors.New("get comment list database error")
		mockCommentsRepo.On("GetCommentList", mock.Anything, filter).Return(nil, repoErr).Once()

		comments, err := commentsService.GetCommentList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Equal(t, repoErr, err)

		mockCommentsRepo.AssertExpectations(t)
	})
}
