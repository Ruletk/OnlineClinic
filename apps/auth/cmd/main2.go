package main

import (
	"auth/config"
	"auth/internal/api"
	"auth/internal/repository"
	"auth/internal/service"
	"github.com/Ruletk/GoMarketplace/pkg/authorization"
	"github.com/Ruletk/GoMarketplace/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

func main() {
	defaultConfig := config.LoadDefaultConfig()
	logging.InitLogger(logging.LogConfig{
		Level:        "debug",
		EnableCaller: true,
		LoggerName:   "auth",
	})

	logging.Logger.Info("Starting the server")

	r := gin.Default()

	r.Use(logging.GinLogger(logging.Logger), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Cookie"},
		AllowCredentials: true,
	}))

	dialer := gomail.NewDialer("smtp.freesmtpservers.com", 25, "", "")
	rbac, err := authorization.NewRBAC(defaultConfig.RBACPath)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to load RBAC")
		panic(err)
	}

	db := ConnectToDB(defaultConfig)

	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Services with no major dependencies
	emailService := service.NewEmailService(dialer)
	jwtService := service.NewJwtService(defaultConfig.Jwt.Algo, defaultConfig.Jwt.Secret)

	// Services with dependencies (cross-service)
	roleService := service.NewRoleService(roleRepo)
	sessionService := service.NewSessionService(sessionRepo)
	authService := service.NewAuthService(authRepo, sessionService, jwtService, emailService)

	authAPI := api.NewAuthAPI(authService, sessionService, roleService)

	router := r.Group("/")
	authAPI.RegisterRoutes(router, rbac)

	err = r.Run(":8080")

	if err != nil {
		return
	}
}

func ConnectToDB(config *config.Config) *gorm.DB {
	dsn := "host=" + config.Database.Host + " user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.Name +
		" port=" + strconv.Itoa(config.Database.Port) +
		" sslmode=disable TimeZone=Asia/Aqtobe"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
