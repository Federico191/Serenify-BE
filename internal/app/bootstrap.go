package app

import (
	"FindIt/internal/auth/delivery"
	"FindIt/internal/auth/repository"
	"FindIt/internal/auth/usecase"
	"FindIt/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	emailPkg "FindIt/pkg/email"
	cronPkg "FindIt/pkg/gocron"
	jwtPkg "FindIt/pkg/jwt"
	"FindIt/pkg/supabase"

	supabaseStorage "github.com/supabase-community/storage-go"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"
)

type BootstrapConfig struct {
	db    *sqlx.DB
	route *gin.Engine
}

func NewBootstrapConfig(db *sqlx.DB, route *gin.Engine) *BootstrapConfig {
	return &BootstrapConfig{db: db, route: route}
}

func (config *BootstrapConfig) Init() {

	// init gocron
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(jakartaTime))
	if err != nil {
		log.Fatalf("cannot initialize gocron: %v", err)
	}
	defer func() {
		_ = scheduler.Shutdown()
	}()

	cron := cronPkg.NewCron(scheduler)

	// init email
	email := emailPkg.NewEmail()

	// init jwt
	jwt := jwtPkg.NewJWT()

	// init supabase
	client := supabaseStorage.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), 
		map[string]string{
			"apikey": os.Getenv("SUPABASE_KEY"),
		})
	
	supabase := supabase.NewSupabaseStorage(client)

	// init repostiory
	authRepo := repository.NewAuthRepo(config.db)

	// init usecase
	authUC := usecase.NewAuthUC(authRepo, email, cron, jwt, supabase)

	// init handler
	authHandler := delivery.NewAuthHandler(authUC)

	// init middleware
	mdlwr := middleware.NewMiddleware(jwt, authUC)

	// start scheduler
	scheduler.Start()

	err = cron.DeleteVerificationCode(func() error {
		return authUC.DeleteVerificationCode()
	})
	if err != nil {
		log.Fatalf("cannot start scheduler: %v", err)
	}

	v1 := config.route.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	// auth routes
	auth := v1.Group("/auth")
	delivery.AuthRoutes(auth, authHandler, mdlwr)

	// start server
	if err := config.route.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
