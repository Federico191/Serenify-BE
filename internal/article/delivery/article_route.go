package delivery

import (
	"github.com/gin-gonic/gin"
)

func ArticleRoutes(articleGroup *gin.RouterGroup, handler *ArticleHandler) {
    articleGroup.GET("/:articleId", handler.GetArticleById)
    articleGroup.GET("", handler.GetAllArticles)
}