package usecase

import (
	"FindIt/internal/auth/repository"
	"FindIt/internal/entity"
	"FindIt/model"
	"FindIt/pkg/email"
	"FindIt/pkg/encode"
	customError "FindIt/pkg/error"
	"FindIt/pkg/gocron"
	"FindIt/pkg/helper"
	jwtPkg "FindIt/pkg/jwt"
	"FindIt/pkg/supabase"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthUCItf interface {
	Register(req model.CreateUserReq) (*model.UserResponse, error)
	Login(req model.LoginUserReq) (string, error)
	GetUserById(id uuid.UUID) (*model.UserResponse, error)
	GetUserByVerificationCode(code sql.NullString) (*entity.User, error)
	RequestResetPassword(email string) error
	ResetPassword(req model.ResetPasswordReq) error
	UpdateUser(req model.UpdateUserReq, id uuid.UUID) error
	DeleteVerificationCode() error
}

type AuthUC struct {
	userRepo repository.AuthRepoItf
	email    email.EmailItf
	cron     gocron.CronItf
	jwt      jwtPkg.JWTItf
	supabase supabase.SupabaseStorageItf
}

func NewAuthUC(userRepo repository.AuthRepoItf, email email.EmailItf,
	cron gocron.CronItf, jwt jwtPkg.JWTItf, supabase supabase.SupabaseStorageItf) AuthUCItf {
	return &AuthUC{userRepo: userRepo, email: email, cron: cron, jwt: jwt, supabase: supabase}
}

// Create implements AuthUCItf.
func (u *AuthUC) Register(req model.CreateUserReq) (*model.UserResponse, error) {
	hashPwd, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	code := helper.GenerateCode()

	parseBirthDate, err := parseDate(req.BirthDate)
	if err != nil {
		return nil, err
	}

	parseCode := sql.NullString{String: code, Valid: true}

	user := &entity.User{
		ID:               uuid.New(),
		FullName:         req.FullName,
		Email:            req.Email,
		Password:         hashPwd,
		BirthDate:        parseBirthDate,
		VerificationCode: parseCode,
	}

	encodedCode := encode.Encode(code)

	err = u.email.SendEmailVerification(user, encodedCode)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.Create(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, customError.ErrEmailAlreadyExists
		}
		return nil, err
	}

	response := convertToUserResponse(user)

	return response, nil
}

// Login implements AuthUCItf.
func (u *AuthUC) Login(req model.LoginUserReq) (string, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", customError.ErrRecordNotFound
		}
		return "", err
	}

	err = helper.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return "", customError.ErrInvalidEmailPassword
	}

	if !user.IsVerified {
		return "", customError.ErrEmailNotVerified
	}

	token, err := u.jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserById implements AuthUCItf.
func (u *AuthUC) GetUserById(id uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrRecordNotFound
		}
		return nil, err
	}

	response := convertToUserResponse(user)

	return response, nil
}

// GetUserByVerificationCode implements AuthUCItf.
func (u *AuthUC) GetUserByVerificationCode(code sql.NullString) (*entity.User, error) {
	user, err := u.userRepo.GetByVerificationCode(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

// RequestResetPassword implements AuthUCItf.
func (u *AuthUC) RequestResetPassword(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customError.ErrRecordNotFound
		}
		return err
	}

	token := uuid.New()

	resetToken:= entity.ResetPasswordToken{
		Token: token,
		UserId: user.ID,
		ExpiredAt: time.Now().Add(time.Hour * 1),
	}

	err = u.userRepo.CreateTokenReset(&resetToken)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return customError.ErrEmailAlreadyExists
		}
		return err
	}

	err = u.email.SendEmailResetPassword(user, token.String())
	if err != nil {
		return nil
	}
	return nil
}

// ResetPassword implements AuthUCItf.
func (u *AuthUC) ResetPassword(req model.ResetPasswordReq) error {
	tokenReset, err := u.userRepo.GetTokenReset(req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customError.ErrRecordNotFound
		}
		return err
	}

	if time.Now().After(tokenReset.ExpiredAt) {
		return errors.New("token expired")
	}

	user, err := u.userRepo.GetById(tokenReset.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customError.ErrRecordNotFound
		}
		return err
	}

	hashPwd, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashPwd

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil

}

// UpdateUser implements AuthUCItf.
func (u *AuthUC) UpdateUser(req model.UpdateUserReq, id uuid.UUID) error {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customError.ErrRecordNotFound
		}
		return err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		hashPwd, err := helper.HashPassword(req.Password)
		if err != nil {
			return err
		}

		user.Password = hashPwd
	}

	if !req.BirthDate.IsZero() {
		user.BirthDate = req.BirthDate
	}

	if req.IsVerified {
		user.IsVerified = req.IsVerified
	}

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUC) DeleteVerificationCode() error {
	expiredCodes, err := u.userRepo.GetExpiredVerificationCode()
	if err != nil {
		return fmt.Errorf("failed to get expired verification code: %v", err)
	}

	if len(expiredCodes) == 0 {
		return fmt.Errorf("no expired verification code found")
	}

	for _, user := range expiredCodes {
		err = u.userRepo.DeleteVerificationCode(user.Email)
		if err != nil {
			continue
		}
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

func parseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
