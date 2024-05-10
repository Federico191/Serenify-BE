package repository

import (
	"FindIt/internal/entity"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepoItf interface {
	GetUserById(id uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepoItf {
	return &UserRepo{db: db}
}

// GetUserById implements UserRepoItf.
func (u *UserRepo) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user entity.User

	err := u.db.QueryRowx(GetUserByIdQuery, id).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) UpdateUser(user *entity.User) error {
	result, err := ur.db.Exec(UpdateUserQuery, user.FullName, user.Email,
		user.Password, user.BirthDate, user.PhotoLink, user.ScoreTest, user.ID)
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
