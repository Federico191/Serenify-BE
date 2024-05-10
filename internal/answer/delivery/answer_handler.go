package delivery

import (
	answerUC "FindIt/internal/answer/usecase"
	"FindIt/model"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnswerHandler struct {
    answerUC answerUC.AnswerUCItf
}

func NewAnswerHandler(answerUC answerUC.AnswerUCItf) *AnswerHandler {
    return &AnswerHandler{
        answerUC: answerUC,
    }
}

func (h *AnswerHandler) EvaluateAnswer(ctx *gin.Context) {
    var req model.AnswerRequest

    userIdCtx, ok := ctx.Get("userId")
	if !ok {
		response.Error(ctx, http.StatusBadRequest, "failed to get user id from context", errors.New(""))
		return 
	}

	userId, ok := userIdCtx.(uuid.UUID)
	if !ok {
		response.Error(ctx, http.StatusBadRequest, "failed to convert user id to uuid", errors.New(""))
		return 
	}

    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to bind json", err)
        return
    }

    result, err := h.answerUC.EvaluateAnswer(req.Answer, userId)
    if err != nil {
        if errors.Is(err, errors.New("answer cannot be empty")) {
            response.Error(ctx, http.StatusBadRequest, "answer cannot be empty", err)
            return
        }
        response.Error(ctx, http.StatusInternalServerError, "failed to evaluate answer", err)
        return
    }

    response.Success(ctx, http.StatusOK, "answer evaluated successfully", result)
}