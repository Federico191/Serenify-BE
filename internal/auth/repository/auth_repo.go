package repository

import (
	"database/sql"
	"fmt"

	"FindIt/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepoItf interface {
	Create(user *entity.User) error
	CreateTokenReset(token *entity.ResetPasswordToken) error
	GetById(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByVerificationCode(code sql.NullString) (*entity.User, error)
	GetTokenReset(token string) (*entity.ResetPasswordToken, error)
	GetExpiredVerificationCode() ([]*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteVerificationCode(email string) error
	DeleteTokenReset(token string) error
}

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) AuthRepoItf {
	return &AuthRepo{db: db}
}

// Create implements AuthRepoItf.
func (ur *AuthRepo) Create(user *entity.User) error {
	result, err := ur.db.Exec(createUserQuery, user.ID, user.FullName, user.BirthDate, user.Email,
		user.Password, user.VerificationCode)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

// CreateTokenReset implements AuthRepoItf.
func (ur *AuthRepo) CreateTokenReset(token *entity.ResetPasswordToken) error {
	result, err := ur.db.Exec(createTokenResetPasswordQuery, token.Token, token.UserId, token.ExpiredAt)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

// GetById implements AuthRepoItf.
func (ur *AuthRepo) GetById(id uuid.UUID) (*entity.User, error) {
	var user entity.User

	err := ur.db.QueryRowx(getUserByIdQuery, id).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail implements AuthRepoItf.
func (ur *AuthRepo) GetByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := ur.db.QueryRowx(getUserByEmailQuery, email).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetTokenReset implements AuthRepoItf.
func (ur *AuthRepo) GetTokenReset(token string) (*entity.ResetPasswordToken, error) {
	var resetPasswordToken entity.ResetPasswordToken

	err := ur.db.QueryRowx(GetTokenResetQuery, token).StructScan(&resetPasswordToken)
	if err != nil {
		return nil, err
	}

	return &resetPasswordToken, nil

}

// GetByVerificationCode implements AuthRepoItf.
func (ur *AuthRepo) GetByVerificationCode(code sql.NullString) (*entity.User, error) {
	var user entity.User

	err := ur.db.QueryRowx(getUserByVerificationCodeQuery, code).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetExpiredVerificationCode implements AuthRepoItf.
func (ur *AuthRepo) GetExpiredVerificationCode() ([]*entity.User, error) {
	rows, err := ur.db.Queryx(getExpiredVerificationCodeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// UpdateUser implements AuthRepoItf.
func (ur *AuthRepo) UpdateUser(user *entity.User) error {
	result, err := ur.db.Exec(updateUserQuery, user.FullName, user.Email,
		user.Password, user.BirthDate, user.IsVerified, user.PhotoLink, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

func (ur *AuthRepo) DeleteVerificationCode(email string) error {
	result, err := ur.db.Exec(deleteVerificationCodeQuery, email)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

// DeleteTokenReset implements AuthRepoItf.
func (ur *AuthRepo) DeleteTokenReset(token string) error {
	result, err := ur.db.Exec(deleteTokenResetPasswordQuery, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}
