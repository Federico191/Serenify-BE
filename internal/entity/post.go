package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Post struct {
    ID uuid.UUID `db:"id"`
    Content string `db:"content"`
    PhotoLink sql.NullString `db:"photo_link"`
    UserID uuid.UUID `db:"user_id"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type PostWithLikeCount struct {
    Post Post
    LikeCount int
}