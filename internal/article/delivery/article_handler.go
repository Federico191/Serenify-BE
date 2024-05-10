package delivery

import (
	"FindIt/internal/article/usecase"
	customError "FindIt/pkg/error"
	"FindIt/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleHandler struct {
    articleUC usecase.ArticleUCItf
}

func NewArticleHandler(articleUC usecase.ArticleUCItf) *ArticleHandler {
    return &ArticleHandler{articleUC: articleUC}
}

func (h *ArticleHandler) GetAllArticles(ctx *gin.Context) {
    articles, err := h.articleUC.GetAllArticles()
    if err != nil {
        if errors.Is(err, customError.ErrRecordNotFound) {
            response.Error(ctx, http.StatusNotFound, "articles not found", err)
            return
        }
        response.Error(ctx, http.StatusInternalServerError, "failed to get articles" , err)
        return
    }

    response.Success(ctx, http.StatusOK, "articles retrieved successfully", articles)
}

func (h *ArticleHandler) GetArticleById(ctx *gin.Context) {
    articleIdParam := ctx.Param("articleId")

    articleId, err := uuid.Parse(articleIdParam)
    if err != nil {
        response.Error(ctx, http.StatusBadRequest, "failed to convert article id to uuid", err)
        return 
    }

    article, err := h.articleUC.GetArticleById(articleId)
    if err != nil {
        if errors.Is(err, customError.ErrRecordNotFound) {
            response.Error(ctx, http.StatusNotFound, "article not found", err)
            return
        }
        response.Error(ctx, http.StatusInternalServerError, "failed to get article" , err)
        return
    }

    response.Success(ctx, http.StatusOK, "article retrieved successfully", article)
}