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
	TokenResetPassword sql.NullString `db:"token_reset_password"`
	IsVerified bool `db:"is_verified"`
	ScoreTest int `db:"score_test"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ResetPasswordToken struct {
	Token uuid.UUID `db:"token"`
	UserId uuid.UUID `db:"user_id"`
	ExpiredAt time.Time `db:"expired_at"`
}