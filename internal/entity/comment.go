package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
    ID int `db:"id"`
    PostID uuid.UUID `db:"post_id"`
    UserID uuid.UUID `db:"user_id"`
    Comment string `db:"comment"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}