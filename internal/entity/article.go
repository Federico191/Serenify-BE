package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Article struct {
    ID uuid.UUID `db:"id"`
    Title string `db:"title"`
    Content string `db:"content"`
    PhotoLink sql.NullString `db:"photo_link"`
    CreatedAt time.Time `db:"created_at"`
}