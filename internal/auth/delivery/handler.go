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

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with create user request
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.CreateUserReq true "Create User Request"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/register [post]
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

	response.Success(ctx, http.StatusCreated, "successfully resgister", user)
}

// Login godoc
// @Summary      Login a user
// @Description  Login a user with login user request
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginUserReq true "Login User Request"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/login [post]
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

// VerifyEmail godoc
// @Summary      Verify email
// @Description  Verify email with verification code
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/verify-email/{verificationCode} [get]
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

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get current user with user id from context
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/current-user [get]
func (uh *AuthHandler) GetCurrentUser(ctx *gin.Context) {
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

// UploadPhoto godoc
// @Summary      Upload photo
// @Description  Upload photo with user id from context and form file
// @Tags         Auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        photo formData file true "Photo"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/upload-photo [post]
func (uh *AuthHandler) UploadPhoto(ctx *gin.Context) {
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

	err = uh.auth.UploadPhoto(req, userId)
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