package usecase

import (
	"FindIt/internal/entity"
	userRepo "FindIt/internal/user/repository"
	"FindIt/model"
	customError "FindIt/pkg/error"
	supabaseStorage "FindIt/pkg/supabase"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type UserUCItf interface {
	GetUserById(id uuid.UUID) (*model.UserResponse, error)
	GetScoreTest(id uuid.UUID) (*model.AnswerResponse, error)
	UploadPhoto(req *model.UploadPhotoReq, id uuid.UUID) error
}

type UserUC struct {
	userRepo userRepo.UserRepoItf
	supabase supabaseStorage.SupabaseStorageItf
}

func NewUserUC(userRepo userRepo.UserRepoItf, supabase supabaseStorage.SupabaseStorageItf) UserUCItf {
	return &UserUC{
		userRepo: userRepo,
		supabase: supabase,
	}
}

// GetUserById implements UserUCItf.
func (u *UserUC) GetUserById(id uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return convertToUserResponse(user), nil
}

// GetScoreTest implements UserUCItf.
func (u *UserUC) GetScoreTest(id uuid.UUID) (*model.AnswerResponse, error) {
	user, err := u.userRepo.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrRecordNotFound
		}
		return nil, err
	}

	return &model.AnswerResponse{
		Score: user.ScoreTest,
	}, nil
}

func (u *UserUC) UploadPhoto(req *model.UploadPhotoReq, id uuid.UUID) error {
	user, err := u.userRepo.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customError.ErrRecordNotFound
		}
		return err
	}

	photoLink, err := u.supabase.Upload(os.Getenv("SUPABASE_BUCKET_USER"), req.Photo)
	if err != nil {
		return fmt.Errorf("failed to upload photo: %w", err)
	}

	user.PhotoLink = sql.NullString{
		String: photoLink,
		Valid:  true,
	}

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func convertToUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		BirthDate:  user.BirthDate,
		PhotoLink:  user.PhotoLink.String,
		IsVerified: user.IsVerified,
	}
}
