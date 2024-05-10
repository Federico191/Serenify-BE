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
    authGroup.POST("/forgot-password", middleware.JwtAuthMiddleware, handler.RequestResetPassword)
    authGroup.GET("/reset-password/:resetToken", handler.ResetPassword)
}