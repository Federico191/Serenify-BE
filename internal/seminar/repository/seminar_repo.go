package repository

import (
	"FindIt/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SeminarRepoItf interface {
	GetAllSeminars() ([]*entity.Seminar, error)
	GetSeminarById(id uuid.UUID) (*entity.Seminar, error)
}

type SeminarRepo struct {
	db *sqlx.DB
}

func NewSeminarRepo(db *sqlx.DB) SeminarRepoItf {
	return &SeminarRepo{db: db}
}

// GetAllSeminars implements SeminarRepoItf.
func (s *SeminarRepo) GetAllSeminars() ([]*entity.Seminar, error) {
	var seminars []*entity.Seminar
    rows, err := s.db.Queryx(getAllSeminarsQuery)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var seminar entity.Seminar
        err = rows.StructScan(&seminar)
        if err != nil {
            return nil, err
        }
        seminars = append(seminars, &seminar)
    }

    return seminars, nil
}

// GetSeminarById implements SeminarRepoItf.
func (s *SeminarRepo) GetSeminarById(id uuid.UUID) (*entity.Seminar, error) {
	var seminar entity.Seminar

    err := s.db.QueryRowx(getSeminarByIdQuery, id).StructScan(&seminar)
    if err != nil {
        return nil, err
    }

    return &seminar, nil
}