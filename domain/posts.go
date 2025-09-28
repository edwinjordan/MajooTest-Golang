package domain

import "time"

type Posts struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostsRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Slug    string `json:"slug" `
}

type CreatePostsRequestSwagger struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostsRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Slug    string `json:"slug" validate:"required"`
}

type UpdatePostsRequestSwagger struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type PostsFilter struct {
	Search string `json:"search" query:"search"`
}
