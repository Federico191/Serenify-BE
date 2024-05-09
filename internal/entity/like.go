package entity

import "github.com/google/uuid"

type PostLike struct {
    UserID uuid.UUID `db:"user_id"`
    PostID uuid.UUID `db:"post_id"`
}

type CommentLike struct {
    CommentID int `db:"comment_id"`
    UserID uuid.UUID `db:"user_id"`
}