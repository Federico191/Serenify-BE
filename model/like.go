package model

import "github.com/google/uuid"

type PostLikeReq struct {
    UserID uuid.UUID `json:"user_id"`
    PostID uuid.UUID `json:"post_id"`
}

type CommentLikeReq struct {
    UserID uuid.UUID `json:"user_id"`
    CommentID int `json:"comment_id"`
}