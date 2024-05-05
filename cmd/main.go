package main

import (
	bootstrap "FindIt/internal/app"
	"FindIt/pkg/db/postgres"
	"log"
	"os"

	docs "FindIt/docs"
   	swaggerfiles "github.com/swaggo/files"
   	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Serenify API
// @version 1.0
// @description Serenify documentation API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	envi := os.Getenv("ENV")
	if err != nil && envi == "" {
		log.Fatalf("cannot load env:%v", err)
	}

	db, err := postgres.DBInit()
	if err != nil {
		log.Fatalf("cannot initialize database: %v", err)
	}

	route := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app := bootstrap.NewBootstrapConfig(db, route)

	app.Init()
}
