package delivery

import (
	"FindIt/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AnswerRoutes(answerGroup *gin.RouterGroup, handler *AnswerHandler, middleware middleware.MiddlewareItf) {
	answerGroup.POST("/evaluate", middleware.JwtAuthMiddleware, handler.EvaluateAnswer)
}
