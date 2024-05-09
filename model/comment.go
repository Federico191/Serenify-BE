package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCommentReq struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	PostID uuid.UUID `json:"post_id" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

type UpdateCommentReq struct {
    ID int `json:"id" binding:"required"`
    UserID uuid.UUID `json:"user_id" binding:"required"`
    Comment string `json:"comment" binding:"required"`
}

type CommentResp struct {
	ID int `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	UserName string `json:"user_name"`
	UserPhoto string `json:"user_photo"`
	PostID uuid.UUID `json:"post_id"`
	Comment string `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
    TotalLikes int `json:"total_likes"`
}