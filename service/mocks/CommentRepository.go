package mocks

import (
	"context"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func NewCommentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentRepository {
	mock := &CommentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

type CommentRepository struct {
	mock.Mock
}

type CommentRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *CommentRepository) EXPECT() *CommentRepository_Expecter {
	return &CommentRepository_Expecter{mock: &_m.Mock}
}

// CreateComment provides a mock function for the type CommentRepository
func (_mock *CommentRepository) CreateComment(ctx context.Context, comment *domain.CreateCommentRequest) (*domain.Comment, error) {
	ret := _mock.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for CreateComment")
	}

	var r0 *domain.Comment
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *domain.CreateCommentRequest) (*domain.Comment, error)); ok {
		return returnFunc(ctx, comment)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *domain.CreateCommentRequest) *domain.Comment); ok {
		r0 = returnFunc(ctx, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Comment)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *domain.CreateCommentRequest) error); ok {
		r1 = returnFunc(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// CommentRepository_CreateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateComment'
type CommentRepository_CreateComment_Call struct {
	*mock.Call
}

// CreateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comment *domain.CreateCommentRequest
func (_e *CommentRepository_Expecter) CreateComment(ctx interface{}, comment interface{}) *CommentRepository_CreateComment_Call {
	return &CommentRepository_CreateComment_Call{Call: _e.mock.On("CreateComment", ctx, comment)}
}

func (_c *CommentRepository_CreateComment_Call) Run(run func(ctx context.Context, comment *domain.CreateCommentRequest)) *CommentRepository_CreateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *domain.CreateCommentRequest
		if args[1] != nil {
			arg1 = args[1].(*domain.CreateCommentRequest)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CommentRepository_CreateComment_Call) Return(comment1 *domain.Comment, err error) *CommentRepository_CreateComment_Call {
	_c.Call.Return(comment1, err)
	return _c
}

func (_c *CommentRepository_CreateComment_Call) RunAndReturn(run func(ctx context.Context, comment *domain.CreateCommentRequest) (*domain.Comment, error)) *CommentRepository_CreateComment_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteComment provides a mock function for the type CommentRepository
func (_mock *CommentRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteComment")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = returnFunc(ctx, id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// CommentRepository_DeleteComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteComment'
type CommentRepository_DeleteComment_Call struct {
	*mock.Call
}

// DeleteComment is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *CommentRepository_Expecter) DeleteComment(ctx interface{}, id interface{}) *CommentRepository_DeleteComment_Call {
	return &CommentRepository_DeleteComment_Call{Call: _e.mock.On("DeleteComment", ctx, id)}
}

func (_c *CommentRepository_DeleteComment_Call) Run(run func(ctx context.Context, id uuid.UUID)) *CommentRepository_DeleteComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 uuid.UUID
		if args[1] != nil {
			arg1 = args[1].(uuid.UUID)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CommentRepository_DeleteComment_Call) Return(err error) *CommentRepository_DeleteComment_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *CommentRepository_DeleteComment_Call) RunAndReturn(run func(ctx context.Context, id uuid.UUID) error) *CommentRepository_DeleteComment_Call {
	_c.Call.Return(run)
	return _c
}

// GetComment provides a mock function for the type UserRepository
func (_mock *CommentRepository) GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetComment")
	}

	var r0 *domain.Comment
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*domain.Comment, error)); ok {
		return returnFunc(ctx, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.Comment); ok {
		r0 = returnFunc(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Comment)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = returnFunc(ctx, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// CommentRepository_GetComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetComment'
type CommentRepository_GetComment_Call struct {
	*mock.Call
}

// GetComment is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *CommentRepository_Expecter) GetComment(ctx interface{}, id interface{}) *CommentRepository_GetComment_Call {
	return &CommentRepository_GetComment_Call{Call: _e.mock.On("GetComment", ctx, id)}
}

func (_c *CommentRepository_GetComment_Call) Run(run func(ctx context.Context, id uuid.UUID)) *CommentRepository_GetComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 uuid.UUID
		if args[1] != nil {
			arg1 = args[1].(uuid.UUID)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CommentRepository_GetComment_Call) Return(comment *domain.Comment, err error) *CommentRepository_GetComment_Call {
	_c.Call.Return(comment, err)
	return _c
}

func (_c *CommentRepository_GetComment_Call) RunAndReturn(run func(ctx context.Context, id uuid.UUID) (*domain.Comment, error)) *CommentRepository_GetComment_Call {
	_c.Call.Return(run)
	return _c
}

// GetCommentList provides a mock function for the type CommentRepository
func (_mock *CommentRepository) GetCommentList(ctx context.Context, filter *domain.CommentFilter) ([]domain.Comment, error) {
	ret := _mock.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentList")
	}

	var r0 []domain.Comment
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *domain.CommentFilter) ([]domain.Comment, error)); ok {
		return returnFunc(ctx, filter)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *domain.CommentFilter) []domain.Comment); ok {
		r0 = returnFunc(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Comment)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *domain.CommentFilter) error); ok {
		r1 = returnFunc(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// CommentRepository_GetCommentList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCommentList'
type CommentRepository_GetCommentList_Call struct {
	*mock.Call
}

// GetCommentList is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *domain.CommentFilter
func (_e *CommentRepository_Expecter) GetCommentList(ctx interface{}, filter interface{}) *CommentRepository_GetCommentList_Call {
	return &CommentRepository_GetCommentList_Call{Call: _e.mock.On("GetCommentList", ctx, filter)}
}

func (_c *CommentRepository_GetCommentList_Call) Run(run func(ctx context.Context, filter *domain.CommentFilter)) *CommentRepository_GetCommentList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *domain.CommentFilter
		if args[1] != nil {
			arg1 = args[1].(*domain.CommentFilter)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CommentRepository_GetCommentList_Call) Return(comments []domain.Comment, err error) *CommentRepository_GetCommentList_Call {
	_c.Call.Return(comments, err)
	return _c
}

func (_c *CommentRepository_GetCommentList_Call) RunAndReturn(run func(ctx context.Context, filter *domain.CommentFilter) ([]domain.Comment, error)) *CommentRepository_GetCommentList_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateComment provides a mock function for the type CommentRepository
func (_mock *CommentRepository) UpdateComment(ctx context.Context, id uuid.UUID, comment *domain.Comment) (*domain.Comment, error) {
	ret := _mock.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComment")
	}

	var r0 *domain.Comment
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, uuid.UUID, *domain.Comment) (*domain.Comment, error)); ok {
		return returnFunc(ctx, id, comment)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, uuid.UUID, *domain.Comment) *domain.Comment); ok {
		r0 = returnFunc(ctx, id, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Comment)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, uuid.UUID, *domain.Comment) error); ok {
		r1 = returnFunc(ctx, id, comment)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// CommentRepository_UpdateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComment'
type CommentRepository_UpdateComment_Call struct {
	*mock.Call
}

// UpdateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - comment *domain.Comment
func (_e *CommentRepository_Expecter) UpdateComment(ctx interface{}, id interface{}, comment interface{}) *CommentRepository_UpdateComment_Call {
	return &CommentRepository_UpdateComment_Call{Call: _e.mock.On("UpdateComment", ctx, id, comment)}
}

func (_c *CommentRepository_UpdateComment_Call) Run(run func(ctx context.Context, id uuid.UUID, comment *domain.Comment)) *CommentRepository_UpdateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 uuid.UUID
		if args[1] != nil {
			arg1 = args[1].(uuid.UUID)
		}
		var arg2 *domain.Comment
		if args[2] != nil {
			arg2 = args[2].(*domain.Comment)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *CommentRepository_UpdateComment_Call) Return(comment1 *domain.Comment, err error) *CommentRepository_UpdateComment_Call {
	_c.Call.Return(comment1, err)
	return _c
}

func (_c *CommentRepository_UpdateComment_Call) RunAndReturn(run func(ctx context.Context, id uuid.UUID, comment *domain.Comment) (*domain.Comment, error)) *CommentRepository_UpdateComment_Call {
	_c.Call.Return(run)
	return _c
}
