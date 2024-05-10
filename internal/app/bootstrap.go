package app

import (
	authDelivery "FindIt/internal/auth/delivery"
	authRepo "FindIt/internal/auth/repository"
	authUC "FindIt/internal/auth/usecase"
	postDelivery "FindIt/internal/post/delivery"
	postRepo "FindIt/internal/post/repository"
	postUC "FindIt/internal/post/usecase"
	userRepo "FindIt/internal/user/repository"
	userUC "FindIt/internal/user/usecase"
	userDelivery "FindIt/internal/user/delivery"
	likeDelivery "FindIt/internal/like/delivery"
	likeRepo "FindIt/internal/like/repository"
	likeUC "FindIt/internal/like/usecase"
	commentDelivery "FindIt/internal/comment/delivery"
	commentRepo "FindIt/internal/comment/repository"
	commentUC "FindIt/internal/comment/usecase"
	answerUC "FindIt/internal/answer/usecase"
	answerDelivery "FindIt/internal/answer/delivery"
	articleUC "FindIt/internal/article/usecase"
	articleDelivery "FindIt/internal/article/delivery"
	articleRepo "FindIt/internal/article/repository"
	seminarDelivery "FindIt/internal/seminar/delivery"
	seminarRepo "FindIt/internal/seminar/repository"
	seminarUC "FindIt/internal/seminar/usecase"
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
	authRepo := authRepo.NewAuthRepo(config.db)
	postRepo := postRepo.NewPostRepo(config.db)
	userRepo := userRepo.NewUserRepo(config.db)
	likeRepo := likeRepo.NewLikeRepo(config.db)
	commentRepo := commentRepo.NewCommentRepo(config.db)
	articleRepo := articleRepo.NewArticleRepo(config.db)
	seminarRepo := seminarRepo.NewSeminarRepo(config.db)

	// init usecase
	authUC := authUC.NewAuthUC(authRepo, email, cron, jwt, supabase)
	postUC := postUC.NewPostUC(postRepo, userRepo, likeRepo, commentRepo, supabase)
	userUC := userUC.NewUserUC(userRepo, supabase)
	likeUC := likeUC.NewLikeUC(likeRepo)
	commentUC := commentUC.NewCommentUC(commentRepo, userRepo, likeRepo)
	answerUC := answerUC.NewAnswerUC(userRepo)
	articleUC := articleUC.NewArticleUC(articleRepo)
	seminarUC := seminarUC.NewSeminarUC(seminarRepo)

	// init handler
	authHandler := authDelivery.NewAuthHandler(authUC)
	userHandler := userDelivery.NewUserHandler(userUC, answerUC)
	postHandler := postDelivery.NewPostHandler(postUC, userUC)
	likeHandler:= likeDelivery.NewLikeHandler(likeUC)
	commentHandler := commentDelivery.NewCommentHandler(commentUC)
	answerHandler := answerDelivery.NewAnswerHandler(answerUC)
	articleHandler := articleDelivery.NewArticleHandler(articleUC)
	seminarHandler := seminarDelivery.NewSeminarHandler(seminarUC)

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
	authDelivery.AuthRoutes(auth, authHandler, mdlwr)

	// user routes
	user := v1.Group("/user")
	userDelivery.UserRoutes(user, userHandler, mdlwr)

	// post routes
	post := v1.Group("/post")
	postDelivery.PostRoutes(post, postHandler, mdlwr)

	// like routes
	likeDelivery.LikeRoutes(post, likeHandler, mdlwr)

	// comment routes
	commentDelivery.CommentRoutes(post, commentHandler, mdlwr)

	// answer routes
	answer := v1.Group("/answer")
	answerDelivery.AnswerRoutes(answer, answerHandler, mdlwr)

	// article routes
	article := v1.Group("/article")
	articleDelivery.ArticleRoutes(article, articleHandler)

	// seminar routes
	seminar := v1.Group("/seminar")
	seminarDelivery.SeminarRoutes(seminar, seminarHandler)

	// start server
	if err := config.route.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
