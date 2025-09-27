package service

import (
	"context"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/google/uuid"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.CreateCommentRequest) (*domain.Comment, error)
	GetCommentList(ctx context.Context, filter *domain.CommentFilter) ([]domain.Comment, error)
	GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	UpdateComment(ctx context.Context, id uuid.UUID, comment *domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

type CommentService struct {
	commentsRepo CommentRepository
}

func NewCommentService(n CommentRepository) *CommentService {
	return &CommentService{
		commentsRepo: n,
	}
}

func (ns *CommentService) CreateComment(
	ctx context.Context,
	u *domain.CreateCommentRequest,
) (*domain.Comment, error) {
	createdComment, err := ns.commentsRepo.CreateComment(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdComment, nil
}

func (us *CommentService) GetComment(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Comment, error) {
	comment, err := us.commentsRepo.GetComment(ctx, id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (us *CommentService) GetCommentList(
	ctx context.Context,
	filter *domain.CommentFilter,
) ([]domain.Comment, error) {
	comments, err := us.commentsRepo.GetCommentList(ctx, filter)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (us *CommentService) UpdateComment(
	ctx context.Context,
	id uuid.UUID,
	u *domain.Comment,
) (*domain.Comment, error) {

	existing, err := us.commentsRepo.GetComment(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	existing.PostID = u.PostID
	existing.UserID = u.UserID
	existing.Body = u.Body

	_, err = us.commentsRepo.UpdateComment(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (us *CommentService) DeleteComment(
	ctx context.Context,
	id uuid.UUID,
) error {

	comment, err := us.commentsRepo.GetComment(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return domain.ErrUserNotFound
	}

	err = us.commentsRepo.DeleteComment(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
