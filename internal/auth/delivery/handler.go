package delivery

import (
	"FindIt/internal/auth/usecase"
	"FindIt/model"
	"FindIt/pkg/encode"
	customError "FindIt/pkg/error"
	"FindIt/pkg/response"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	auth usecase.AuthUCItf
}

func NewAuthHandler(auth usecase.AuthUCItf) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (uh *AuthHandler) Register(ctx *gin.Context) {
	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	user, err := uh.auth.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrEmailAlreadyExists):
			response.Error(ctx, http.StatusConflict, "email already exists", err)
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to register user", err)
		}
		return
	}

	response.Success(ctx, http.StatusCreated, "successfully register", user)
}

func (uh *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	token, err := uh.auth.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "user not found", err)
		case errors.Is(err, customError.ErrInvalidEmailPassword):
			response.Error(ctx, http.StatusUnauthorized, "invalid email or password", err)
		case errors.Is(err, customError.ErrEmailNotVerified):
			response.Error(ctx, http.StatusForbidden, "email not verified", err)
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to login user", err)
		}
		return
	}

	response.Success(ctx, http.StatusOK, "successfully logged in", token)
}

func (uh *AuthHandler) VerifyEmail(ctx *gin.Context) {
	codeParam := ctx.Param("verificationCode")
	log.Println(codeParam)
	

	codeStr, err := encode.Decode(codeParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to decode verification code", err)
		return
	}
	
	code := sql.NullString{String: codeStr, Valid: true}

	user, err := uh.auth.GetUserByVerificationCode(code)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to get user by verification code", err)
		return
	}



	if user.VerificationCode != code {
		response.Error(ctx, http.StatusBadRequest, "invalid verification code", err)
		return
	}

	err = uh.auth.UpdateUser(model.UpdateUserReq{IsVerified: true}, user.ID)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	response.Success(ctx, http.StatusOK, "email verified", user.IsVerified)
}

func (uh *AuthHandler) GetCurrentUser(ctx *gin.Context) {
	userId := getUserId(ctx)

	user, err := uh.auth.GetUserById(userId)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to get user by id", err)
		return 
	}

	response.Success(ctx, http.StatusOK, "success get current user", user)
}

func (uh *AuthHandler) RequestResetPassword(ctx *gin.Context) {
	var req model.RequestResetPassword

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	err := uh.auth.RequestResetPassword(req.Email)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "user not found", err)
		} else if errors.Is(err, customError.ErrEmailNotVerified) {
			response.Error(ctx, http.StatusForbidden, "email not verified", err)
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to request reset password", err)
		return
	}
}

func (uh *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req model.ResetPasswordReq

	tokenParam := ctx.Param("resetToken")

	req.Token = tokenParam

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	err := uh.auth.ResetPassword(req)
	if err != nil {

		response.Error(ctx, http.StatusInternalServerError, "failed to reset password", err)
	}
}

func getUserId(ctx *gin.Context) uuid.UUID {
	userIdCtx, ok := ctx.Get("userId")
	if !ok {
		response.Error(ctx, http.StatusBadRequest, "failed to get user id from context", errors.New(""))
		return uuid.UUID{}
	}

	userId, ok := userIdCtx.(uuid.UUID)
	if !ok {
		response.Error(ctx, http.StatusBadRequest, "failed to convert user id to uuid", errors.New(""))
		return uuid.UUID{}
	}

	return userId
}