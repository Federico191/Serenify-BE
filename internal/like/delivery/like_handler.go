package delivery

import (
	likeUC "FindIt/internal/like/usecase"
	customError "FindIt/pkg/error"
	"strconv"

	"FindIt/model"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeHandler struct {
    likeUC likeUC.LikeUCItf
}

func NewLikeHandler(likeUC likeUC.LikeUCItf) *LikeHandler {
    return &LikeHandler{likeUC: likeUC}
}

func (h *LikeHandler) CreatePostLike(ctx *gin.Context) {
    userId := getUserId(ctx)

	postId := getPostId(ctx)

	req := model.PostLikeReq{
		UserID: userId,
		PostID: postId,
	}

	err := h.likeUC.CreatePostLike(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordAlreadyExists):
			response.Error(ctx, http.StatusConflict, "failed to create like", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to create like", err)
			return
		}
	}

	response.Success(ctx, http.StatusCreated, "post like created successfully", nil)
}

func (h *LikeHandler) CreateCommentLike(ctx *gin.Context) {
    userId := getUserId(ctx)

	commentId := getCommentId(ctx)

	req := model.CommentLikeReq{
		UserID: userId,
		CommentID: commentId,
	}

	err := h.likeUC.CreateCommentLike(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordAlreadyExists):
			response.Error(ctx, http.StatusConflict, "failed to create like", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to create like", err)
			return
		}
	}

	response.Success(ctx, http.StatusCreated, "comment like created successfully", nil)
}

func (h *LikeHandler) DeletePostLike(ctx *gin.Context) {
    userId := getUserId(ctx)

	postId := getPostId(ctx)

	req := model.PostLikeReq{
		UserID: userId,
		PostID: postId,
	}

	err := h.likeUC.DeletePostLike(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "failed to delete like", err)
			return
		case errors.Is(err, customError.ErrNotAuthorize):
			response.Error(ctx, http.StatusForbidden, "failed to delete like", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to delete like", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "post like deleted successfully", nil)
}

func (h *LikeHandler) DeleteCommentLike(ctx *gin.Context) {
    userId := getUserId(ctx)

	commentId := getCommentId(ctx)

	req := model.CommentLikeReq{
		UserID: userId,
		CommentID: commentId,
	}

	err := h.likeUC.DeleteCommentLike(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "failed to delete like", err)
			return
		case errors.Is(err, customError.ErrNotAuthorize):
			response.Error(ctx, http.StatusForbidden, "failed to delete like", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to delete like", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "comment like deleted successfully", nil)
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

func getPostId(ctx *gin.Context) uuid.UUID {
	postIdCtx := ctx.Param("postId")

	postId, err := uuid.Parse(postIdCtx)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to convert post id to uuid", err)
		return uuid.UUID{}
	}

	return postId
}

func getCommentId(ctx *gin.Context) int {
	CommentIdCtx := ctx.Param("commentId")

	CommentId, err := strconv.Atoi(CommentIdCtx)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to convert post id to uuid", err)
		return 0
	}

	return CommentId
}