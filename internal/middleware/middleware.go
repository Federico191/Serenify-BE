package middleware

import (
	"FindIt/internal/auth/usecase"
	jwtPkg "FindIt/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type MiddlewareItf interface {
	JwtAuthMiddleware(ctx *gin.Context)
}

type Middleware struct {
	jwt    jwtPkg.JWTItf
	authUC usecase.AuthUCItf
}

func NewMiddleware(jwt jwtPkg.JWTItf, authUC usecase.AuthUCItf) MiddlewareItf {
	return &Middleware{jwt: jwt, authUC: authUC}
}
