package usecase

import (
	"FindIt/internal/entity"
	"FindIt/internal/seminar/repository"
	"FindIt/model"
    customError "FindIt/pkg/error"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type SeminarUCItf interface {
	GetAllSeminars() ([]*model.SeminarResponse, error)
	GetSeminarById(id uuid.UUID) (*model.SeminarDetailResponse, error)
}

type SeminarUC struct {
	seminarRepo repository.SeminarRepoItf
}

func NewSeminarUC(seminarRepo repository.SeminarRepoItf) SeminarUCItf {
	return &SeminarUC{seminarRepo}
}

// GetAllSeminars implements SeminarUCItf.
func (s *SeminarUC) GetAllSeminars() ([]*model.SeminarResponse, error) {
	seminars, err := s.seminarRepo.GetAllSeminars()
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, customError.ErrRecordNotFound
        }
        return nil, err
    }

    var responses []*model.SeminarResponse
    for _, seminar := range seminars {
        response := &entity.Seminar{
            ID:        seminar.ID,
            Title:     seminar.Title,
            Time:      seminar.Time,
            Place:     seminar.Place,
            Price:     seminar.Price,
            PhotoLink: seminar.PhotoLink,
            CreatedAt: seminar.CreatedAt,
        }

        responses = append(responses, &model.SeminarResponse{
            ID:        response.ID,
            Title:     response.Title,
            Time:      response.Time,
            Place:     response.Place,
            Price:     response.Price,
            PhotoLink: response.PhotoLink.String,
        })
    }

    return responses, nil
}

// GetSeminarById implements SeminarUCItf.
func (s *SeminarUC) GetSeminarById(id uuid.UUID) (*model.SeminarDetailResponse, error) {
    seminar, err := s.seminarRepo.GetSeminarById(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, customError.ErrRecordNotFound
        }
        return nil, err
    }

    return &model.SeminarDetailResponse{
        ID:        seminar.ID,
        Title:     seminar.Title,
        Time:      seminar.Time,
        Place:     seminar.Place,
        Price:     seminar.Price,
        PhotoLink: seminar.PhotoLink.String,
        Description: seminar.Description,
    }, nil
}