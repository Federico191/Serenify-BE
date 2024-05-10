package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(postGroup *gin.RouterGroup, handler *PostHandler, middleware middleware.MiddlewareItf) {
    postGroup.POST("/create", middleware.JwtAuthMiddleware, handler.CreatePost)
    postGroup.GET("", handler.GetAllPosts)
    postGroup.GET("/:postId", handler.GetPost)
    postGroup.PATCH("/:postId", middleware.JwtAuthMiddleware, handler.UpdatePost)
    postGroup.DELETE("/:postId", middleware.JwtAuthMiddleware, handler.DeletePost)
}