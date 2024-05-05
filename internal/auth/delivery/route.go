package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(authGroup *gin.RouterGroup, handler *AuthHandler, middleware middleware.MiddlewareItf) {
    authGroup.POST("/register", handler.Register)
    authGroup.POST("/login", handler.Login)
    authGroup.GET("/verify-email/:verificationCode", handler.VerifyEmail)
    authGroup.GET("/current-user", middleware.JwtAuthMiddleware, handler.GetCurrentUser)
    authGroup.POST("/upload-photo", middleware.JwtAuthMiddleware, handler.UploadPhoto)

}