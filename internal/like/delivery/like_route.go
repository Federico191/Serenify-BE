package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(postGroup *gin.RouterGroup, handler *LikeHandler, middleware middleware.MiddlewareItf) {
    postGroup.POST("/:postId/like", middleware.JwtAuthMiddleware, handler.CreatePostLike)
    postGroup.POST("/:postId/comments/:commentId/like", middleware.JwtAuthMiddleware, handler.CreateCommentLike)
    postGroup.DELETE("/:postId/like", middleware.JwtAuthMiddleware, handler.DeletePostLike)
    postGroup.DELETE("/:postId/comments/:commentId/like", middleware.JwtAuthMiddleware, handler.DeleteCommentLike)
}