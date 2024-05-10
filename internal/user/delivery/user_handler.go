package delivery

import (
	userUC "FindIt/internal/user/usecase"
	answerUC "FindIt/internal/answer/usecase"
	"FindIt/model"
	customError "FindIt/pkg/error"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUC userUC.UserUCItf
	answerUC answerUC.AnswerUCItf
}

func NewUserHandler(userUC userUC.UserUCItf, answerUC answerUC.AnswerUCItf) *UserHandler {
	return &UserHandler{
		userUC: userUC,
		answerUC: answerUC,
	}
}

func (uh *UserHandler) UploadPhoto(ctx *gin.Context) {
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

	file, err := ctx.FormFile("photo")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to get file", err)
		return
	}


	req := &model.UploadPhotoReq{
		Photo: file,
	}

	err = uh.userUC.UploadPhoto(req, userId)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to upload photo", err)
		return 
	}

	response.Success(ctx, http.StatusOK, "success upload photo", file.Filename)
}

func (uh *UserHandler) GetScoreTest(ctx *gin.Context) {
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

	result, err := uh.userUC.GetScoreTest(userId)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to get score", err)
		return
	}

	description := uh.answerUC.CheckScore(result.Score)

	result.Description = description 

	response.Success(ctx, http.StatusOK, "success get score", result)
}