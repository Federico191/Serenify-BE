package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(userGroup *gin.RouterGroup, handler *UserHandler, middleware middleware.MiddlewareItf) {
    userGroup.POST("/upload-photo", middleware.JwtAuthMiddleware, handler.UploadPhoto)
    userGroup.GET("/score", middleware.JwtAuthMiddleware, handler.GetScoreTest)
}