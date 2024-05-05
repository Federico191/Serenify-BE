package repository

import (
	"database/sql"
	"fmt"
	"log"

	"FindIt/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepoItf interface {
	Create(user *entity.User) error
	GetById(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByVerificationCode(code sql.NullString) (*entity.User, error)
	GetEmailExist(email string) bool
	GetExpiredVerificationCode() ([]*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteVerificationCode(email string) error
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

// GetByVerificationCode implements AuthRepoItf.
func (ur *AuthRepo) GetByVerificationCode(code sql.NullString) (*entity.User, error) {
	var user entity.User

	err := ur.db.QueryRowx(getUserByVerificationCodeQuery, code).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetEmailExist implements AuthRepoItf.
func (ur *AuthRepo) GetEmailExist(email string) bool {
	var count int

	err := ur.db.QueryRowx(getUserByEmailQuery, email).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
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
	log.Printf("success delete verification code for %s\n", email)

	return nil
}
