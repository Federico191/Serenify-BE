package delivery

import "github.com/gin-gonic/gin"

func SeminarRoutes(seminar *gin.RouterGroup, seminarHandler *SeminarHandler) {
    seminar.GET("", seminarHandler.GetAllSeminars)
    seminar.GET("/:seminarId", seminarHandler.GetSeminarById)
}