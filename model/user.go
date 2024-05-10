package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	BirthDate string `json:"birth_date" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type UserResponse struct {
	ID uuid.UUID `json:"id"`
	FullName string `json:"name"`
	Email string `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	PhotoLink string `json:"photo_link"`
	IsVerified bool `json:"is_verified"`
}

type LoginUserReq struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	BirthDate time.Time `json:"birth_date" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	IsVerified bool `json:"is_verified"`
}

type UploadPhotoReq struct {
	Photo *multipart.FileHeader `form:"photo"`
}

type RequestResetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordReq struct {
	Token string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
}