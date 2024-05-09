package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(postGroup *gin.RouterGroup, handler *CommentHandler, middleware middleware.MiddlewareItf) {
    postGroup.POST("/:postId/comments", middleware.JwtAuthMiddleware, handler.CreateComment)
    postGroup.PATCH("/:postId/comments/:commentId", middleware.JwtAuthMiddleware, handler.UpdateComment)
    postGroup.DELETE("/:postId/comments/:commentId", middleware.JwtAuthMiddleware, handler.DeleteComment)
}