package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Seminar struct {
    ID uuid.UUID `db:"id"`
    Title string `db:"title"`
    Time string `db:"time"`
    Place string `db:"place"`
    Price int `db:"price"`
    Description string `db:"description"`
    PhotoLink sql.NullString `db:"photo_link"`
    CreatedAt time.Time `db:"created_at"`
}