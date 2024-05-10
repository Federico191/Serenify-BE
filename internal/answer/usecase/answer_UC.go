package usecase

import (
	userRepo "FindIt/internal/user/repository"
	"FindIt/model"
	"errors"

	"github.com/google/uuid"
)

type AnswerUCItf interface {
	EvaluateAnswer(answer []string, userId uuid.UUID) (*model.AnswerResponse, error)
    CheckScore(score int) string
}

type AnswerUC struct {
	userRepo userRepo.UserRepoItf
}

var correctAnswers = []string{
	"B", "B", "B", "B", "S",
	"S", "S", "S", "S", "S",
	"S", "S", "S", "B", "S",
	"B", "B", "B", "S", "S",
}

func NewAnswerUC(userRepo userRepo.UserRepoItf) AnswerUCItf {
	return &AnswerUC{userRepo: userRepo}
}

// EvaluateAnswer implements AnswerUCItf.
func (a *AnswerUC) EvaluateAnswer(answer []string, userId uuid.UUID) (*model.AnswerResponse, error) {
	if len(answer) == 0 {
		return nil, errors.New("answer cannot be empty")
	}

    user, err := a.userRepo.GetUserById(userId)
    if err != nil {
        return nil, err
    }

    score := 0

    if user.ScoreTest > 0 {
        user.ScoreTest = score
    }

    for i := 0; i < len(answer); i++ {
        if answer[i] == correctAnswers[i] {
            score += 5
        }
    }

    description := a.CheckScore(score)

    user.ScoreTest = score

    err = a.userRepo.UpdateUser(user)
    if err != nil {
        return nil, err
    }

    return &model.AnswerResponse{Score: score, Description: description}, nil

}

func (a *AnswerUC) CheckScore(score int) string {
    description := "Hasil tes kamu sangat mengesankan! Teruslah menjaga dan memberikan perhatian pada kesejahteraan mental kamu dan orang disekitar kamu."
    if score < 65 {
        description = "Hasil tes kamu menunjukkan perlunya perhatian lebih. Kami siap membantu menuju perjalanan yang lebih baik."
    }
    return description
}