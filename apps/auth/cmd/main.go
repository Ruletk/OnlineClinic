package main

import (
	"auth/internal/api"
	"auth/internal/repository"
	"auth/internal/service"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		panic(err) // The initialization of the app failed, without a config we can't do anything
	}
	logging.InitLogger(*cfg)

	r := gin.Default()

	r.Use(logging.GinLogger(logging.Logger), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Cookie"},
		AllowCredentials: true,
	}))

	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to the database")
		panic(err) // Also, the initialization of the app failed, without a database we can't do anything
	}

	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// TODO: Add jwt settings to the config
	//jwtService := service.NewJwtService(defaultConfig.Jwt.Algo, defaultConfig.Jwt.Secret)
	jwtService := service.NewJwtService(jwt.SigningMethodHS256, "This is a secret key temp for avoid errors")

	roleService := service.NewRoleService(roleRepo)
	sessionService := service.NewSessionService(sessionRepo)
	authService := service.NewAuthService(authRepo, sessionService, jwtService)

	authAPI := api.NewAuthAPI(authService, sessionService, roleService)

	router := r.Group("/")
	authAPI.RegisterRoutes(router)

	err = r.Run(":8080")

	if err != nil {
		return
	}
}
