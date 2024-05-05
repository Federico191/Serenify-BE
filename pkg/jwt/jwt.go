package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTItf interface {
	CreateToken(id uuid.UUID) (string, error)
	VerifyToken(token string) (uuid.UUID, error)
}

type JWT struct {
}

type Claims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

func NewJWT() JWTItf {
    return &JWT{}
}

func (j *JWT) CreateToken(id uuid.UUID) (string, error) {
    expired, err := strconv.Atoi(os.Getenv("EXPIRED_TOKEN"))
    if err != nil {
        return "", fmt.Errorf("failed to convert expired token: %v", err)
    }

    expirationnTime := time.Now().Add(time.Duration(expired) * time.Hour)

    claims := &Claims{
        UserId: id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationnTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %v", err)
    }

    return signedToken, nil
}

func (j *JWT) VerifyToken(tokenString string) (uuid.UUID, error) {
    var claims Claims

    token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("SECRET_TOKEN")), nil
    })
    if err != nil {
        return uuid.UUID{}, fmt.Errorf("failed to parse token: %v", err)
    }

    if !token.Valid {
        return uuid.UUID{}, fmt.Errorf("invalid token")
    }

    return claims.UserId, nil


}