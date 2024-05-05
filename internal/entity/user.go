package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id"`
	FullName string `db:"full_name"`
	Email string `db:"email"`
	Password string `db:"password"`
	BirthDate time.Time `db:"birth_date"`
	PhotoLink sql.NullString `db:"photo_link"`
	VerificationCode sql.NullString `db:"verification_code"`
	IsVerified bool `db:"is_verified"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}