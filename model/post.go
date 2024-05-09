package model

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreatePostReq struct {
	UserID  uuid.UUID             `form:"user_id" binding:"required"`
	Content string                `form:"content" binding:"required"`
	Photo   *multipart.FileHeader `form:"photo"`
}

type UpdatePostReq struct {
	ID      uuid.UUID             `form:"id" binding:"required"`
	UserID  uuid.UUID             `form:"user_id" binding:"required"`
	Content string                `form:"content" binding:"required"`
	Photo   *multipart.FileHeader `form:"photo"`
}

type PostResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserPhoto     string    `json:"user_photo"`
	Content       string    `json:"content"`
	PhotoLink     string    `json:"photo_link"`
	TotalLikes    int       `json:"total_likes"`
	TotalComments int       `json:"total_comments"`
}

type PostDetailResponse struct {
	ID            uuid.UUID           `json:"id"`
	UserID        uuid.UUID           `json:"user_id"`
	UserName      string              `json:"user_name"`
	UserPhoto     string              `json:"user_photo"`
	Content       string              `json:"content"`
	PhotoLink     string              `json:"photo_link"`
	TotalLikes    int                 `json:"total_likes"`
	TotalComments int                 `json:"total_comments"`
	Comments      []CommentResp `json:"comments"`
}
