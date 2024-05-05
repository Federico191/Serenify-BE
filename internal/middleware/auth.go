package middleware

import (
	"FindIt/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware implements MiddlewareItf.
func (m *Middleware) JwtAuthMiddleware(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
    if bearer == "" {
        response.Error(ctx, http.StatusUnauthorized, "Authorization token is required", nil)
        ctx.Abort()
        return
    }

    token := strings.Split(bearer, " ")[1]
    userId, err := m.jwt.VerifyToken(token)
    if err != nil {
        response.Error(ctx, http.StatusUnauthorized, "Invalid token", err)
        ctx.Abort()
        return 
    }

    ctx.Set("userId", userId)
    ctx.Next()
}