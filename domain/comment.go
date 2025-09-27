package domain

import "time"

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

type UpdateCommentRequest struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

type CommentFilter struct {
	Search string `json:"search" query:"search"`
}
