package delivery

import (
	postUC "FindIt/internal/post/usecase"
	userUC "FindIt/internal/user/usecase"
	customError "FindIt/pkg/error"
	"strconv"

	"FindIt/model"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postUC postUC.PostUCItf
	userUC userUC.UserUCItf
}

func NewPostHandler(postUC postUC.PostUCItf, userUC userUC.UserUCItf) *PostHandler {
	return &PostHandler{
		postUC: postUC,
		userUC: userUC,
	}
}

func (p *PostHandler) CreatePost(ctx *gin.Context) {
	userId := getUserId(ctx)

	photo, err := ctx.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			req := &model.CreatePostReq{
				UserID: userId,
				Photo:  nil,
			}

			if err := ctx.ShouldBind(&req); err != nil {
				response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
				return
			}

			post, err := p.postUC.CreatePost(req)
			if err != nil {
				switch {
				case errors.Is(err, customError.ErrRecordAlreadyExists):
					response.Error(ctx, http.StatusConflict, "failed to create post", err)
				default:
					response.Error(ctx, http.StatusInternalServerError, "failed to create post", err)
					return
				}
			}
			response.Success(ctx, http.StatusCreated, "post created successfully", post)
			return
		}
		response.Error(ctx, http.StatusBadRequest, "failed to get photo", err)
		return 
	}

	req := &model.CreatePostReq{
		UserID: userId,
		Photo:  photo,
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	post, err := p.postUC.CreatePost(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordAlreadyExists):
			response.Error(ctx, http.StatusConflict, "failed to create post", err)
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to create post", err)
			return
		}
	}

	response.Success(ctx, http.StatusCreated, "post created successfully", post)
}

func (p *PostHandler) GetPost(ctx *gin.Context) {
	postId := getPostId(ctx)

	post, err := p.postUC.GetPostByPostId(postId)
	if err != nil {
		if errors.Is(err, customError.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "failed to get post", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to get post", err)
		return
	}

	response.Success(ctx, http.StatusOK, "post retrieved successfully", post)
}

func (p *PostHandler) GetAllPosts(ctx *gin.Context) {
	pageQuery := ctx.Query("page")
	page , err := strconv.Atoi(pageQuery)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to get posts", err)
		return 
	}
	
	posts, err := p.postUC.GetAllPosts(page)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "failed to get posts", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to get posts", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "posts retrieved successfully", posts)
}

func (p *PostHandler) UpdatePost(ctx *gin.Context) {
	postId := getPostId(ctx)

	userId := getUserId(ctx)

	photo, err := ctx.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			req := &model.UpdatePostReq{
				ID: postId,
				UserID: userId,
				Photo:  nil,
			}

			if err := ctx.ShouldBind(&req); err != nil {
				response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
				return
			}

			post, err := p.postUC.UpdatePost(req)
			if err != nil {
				switch {
				case errors.Is(err, customError.ErrRecordAlreadyExists):
					response.Error(ctx, http.StatusConflict, "failed to create post", err)
				default:
					response.Error(ctx, http.StatusInternalServerError, "failed to create post", err)
					return
				}
			}
			response.Success(ctx, http.StatusCreated, "post created successfully", post)
			return
		}
		response.Error(ctx, http.StatusBadRequest, "failed to get photo", err)
		return 
	}

	req := &model.UpdatePostReq{
		Photo: photo,
		ID: postId,
		UserID: userId,
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	post, err := p.postUC.UpdatePost(req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrRecordNotFound):
			response.Error(ctx, http.StatusNotFound, "failed to update post", err)
			return
		case errors.Is(err, customError.ErrNotAuthorize):
			response.Error(ctx, http.StatusForbidden, "failed to update post", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to update post", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "post updated successfully", post)
}

func (p *PostHandler) DeletePost(ctx *gin.Context) {
	postId := getPostId(ctx)

	userId := getUserId(ctx)

	err := p.postUC.DeletePost(userId, postId)
	if err != nil {
		switch {
        case errors.Is(err, customError.ErrRecordNotFound):
            response.Error(ctx, http.StatusNotFound, "failed to delete post", err)
            return
        case errors.Is(err, customError.ErrNotAuthorize):
            response.Error(ctx, http.StatusUnauthorized, "failed to delete post", err)
            return        
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed to delete post", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "post deleted successfully", nil)
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