package service

import (
	"context"

	"github.com/edwinjordan/MajooTest-Golang/domain"
	"github.com/google/uuid"
)

type PostsRepository interface {
	CreatePosts(ctx context.Context, posts *domain.CreatePostsRequest) (*domain.Posts, error)
	GetPostsList(ctx context.Context, filter *domain.PostsFilter) ([]domain.Posts, error)
	GetPosts(ctx context.Context, id uuid.UUID) (*domain.Posts, error)
	UpdatePosts(ctx context.Context, id uuid.UUID, posts *domain.Posts) (*domain.Posts, error)
	DeletePosts(ctx context.Context, id uuid.UUID) error
}

type PostsService struct {
	postsRepo PostsRepository
}

func NewPostsService(n PostsRepository) *PostsService {
	return &PostsService{
		postsRepo: n,
	}
}

func (ns *PostsService) CreatePosts(
	ctx context.Context,
	u *domain.CreatePostsRequest,
) (*domain.Posts, error) {
	createdPosts, err := ns.postsRepo.CreatePosts(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdPosts, nil
}

func (us *PostsService) GetPosts(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Posts, error) {
	posts, err := us.postsRepo.GetPosts(ctx, id)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (us *PostsService) UpdatePosts(
	ctx context.Context,
	id uuid.UUID,
	u *domain.Posts,
) (*domain.Posts, error) {

	existing, err := us.postsRepo.GetPosts(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	existing.Title = u.Title
	existing.Slug = u.Slug
	existing.Content = u.Content

	_, err = us.postsRepo.UpdatePosts(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (us *PostsService) DeletePosts(
	ctx context.Context,
	id uuid.UUID,
) error {

	posts, err := us.postsRepo.GetPosts(ctx, id)
	if err != nil {
		return err
	}
	if posts == nil {
		return domain.ErrUserNotFound
	}

	err = us.postsRepo.DeletePosts(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *PostsService) GetPostsList(ctx context.Context, filter *domain.PostsFilter) ([]domain.Posts, error) {
	postsList, err := us.postsRepo.GetPostsList(ctx, filter)
	if err != nil {
		return nil, err
	}
	return postsList, nil
}
