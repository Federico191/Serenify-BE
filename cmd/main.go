package main

import (
	bootstrap "FindIt/internal/app"
	"FindIt/pkg/db/postgres"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	// postgres.SeedInit(db)

	route := gin.Default()

	app := bootstrap.NewBootstrapConfig(db, route)

	app.Init()
}
