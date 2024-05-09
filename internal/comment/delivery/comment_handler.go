package delivery

import (
	"FindIt/internal/comment/usecase"
	"FindIt/model"
	customError "FindIt/pkg/error"
	"FindIt/pkg/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
    CommentUC usecase.CommentUCItf
}

func NewCommentHandler(CommentUC usecase.CommentUCItf) *CommentHandler {
    return &CommentHandler{CommentUC: CommentUC}
}

func (c *CommentHandler) CreateComment(ctx *gin.Context) {
    var req model.CreateCommentReq

    userId := getUserId(ctx)

    postId := getPostId(ctx)

    req.UserID = userId

    req.PostID = postId

    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to bind json", err)
        return
    }

    comment, err := c.CommentUC.CreateComment(req)
    if err != nil {
        switch {
		case errors.Is(err, customError.ErrRecordAlreadyExists):
		    response.Error(ctx, http.StatusConflict, "failed to create comment", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to create comment", err)
			return
		}
    }

    response.Success(ctx, http.StatusOK, "comment created successfully", comment)

}

func (c *CommentHandler) UpdateComment(ctx *gin.Context) {
    commentIdCtx := ctx.Param("commentId")

    commentId, err := strconv.Atoi(commentIdCtx)
    if err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to convert comment id to int", err)
        return
    }

    userId := getUserId(ctx)

    req := model.UpdateCommentReq{
        ID: commentId,
        UserID: userId,
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to bind json", err)
        return
    }

    comment, err := c.CommentUC.UpdateComment(req)
    if err != nil {
        switch {
        case errors.Is(err, customError.ErrRecordNotFound):
            response.Error(ctx, http.StatusNotFound, "failed to update comment", err)
            return
        case errors.Is(err, customError.ErrNotAuthorize):
            response.Error(ctx, http.StatusForbidden, "failed to update comment", err)
            return
        default:
            response.Error(ctx, http.StatusInternalServerError, "failed to update comment", err)
            return 
        }

    }

    response.Success(ctx, http.StatusOK, "comment updated successfully", comment)
}

func (c *CommentHandler) DeleteComment(ctx *gin.Context) {
    commentIdCtx := ctx.Param("commentId")

	commentId, err := strconv.Atoi(commentIdCtx)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to convert comment id to int", err)
		return
	}

    userId := getUserId(ctx)

	err = c.CommentUC.DeleteComment(userId, commentId)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "failed to delete comment", err)
			return
		case errors.Is(err, customError.ErrNotAuthorize):
			response.Error(ctx, http.StatusForbidden, "failed to delete comment", err)
            return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to delete comment", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "comment deleted successfully", nil)   
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