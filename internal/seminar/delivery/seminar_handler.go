package delivery

import (
	"FindIt/internal/seminar/usecase"
	customError "FindIt/pkg/error"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeminarHandler struct {
    seminarUC usecase.SeminarUCItf
}

func NewSeminarHandler(seminarUC usecase.SeminarUCItf) *SeminarHandler {
    return &SeminarHandler{
        seminarUC: seminarUC,
    }
}

func (h *SeminarHandler) GetAllSeminars(ctx *gin.Context) {
    seminars, err := h.seminarUC.GetAllSeminars()
    if err != nil {
        if errors.Is(err, customError.ErrRecordNotFound) {
            response.Error(ctx, http.StatusNotFound, "failed to get all seminars", err)
            return
        }
        response.Error(ctx, http.StatusInternalServerError, "failed to get all seminars", err)
        return
    }

    response.Success(ctx, http.StatusOK, "all seminars retrieved successfully", seminars)
}

func (h *SeminarHandler) GetSeminarById(ctx *gin.Context) {
    seminarIdParam := ctx.Param("seminarId")

    seminarId, err := uuid.Parse(seminarIdParam)
    if err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to convert seminar id to uuid", err)
        return 
    }

    seminar, err := h.seminarUC.GetSeminarById(seminarId)
    if err != nil {
        if errors.Is(err, customError.ErrRecordNotFound) {
            response.Error(ctx, http.StatusNotFound, "failed to get seminar by id", err)
            return
        }
        response.Error(ctx, http.StatusInternalServerError, "failed to get seminar by id", err)
        return
    }

    response.Success(ctx, http.StatusOK, "seminar retrieved successfully", seminar)
}