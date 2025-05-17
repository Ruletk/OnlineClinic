package main

import (
	"auth/internal/api"
	nats2 "auth/internal/nats"
	"auth/internal/repository"
	"auth/internal/service"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"strconv"
)

func main() {
	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		panic(err) // The initialization of the app failed, without a config we can't do anything
	}
	logging.InitLogger(*cfg)

	logging.Logger.Debug("Creating Gin router")
	r := gin.Default()

	logging.Logger.Debug("Setting up CORS middleware")
	r.Use(logging.GinLogger(logging.Logger), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Cookie"},
		AllowCredentials: true,
	}))

	logging.Logger.Debug("Setting up NATS connection")
	natsConn, err := nats.Connect(cfg.Nats.Url)
	_ = natsConn
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to connect to NATS. Disabling NATS features.")
		natsConn = nil
	}

	logging.Logger.Debug("Setting up database connection")
	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to the database")
		panic(err) // Also, the initialization of the app failed, without a database we can't do anything
	}

	logging.Logger.Debug("Starting repositories")
	natsPublisher := nats2.NewPublisher(natsConn)
	logging.Logger.Debugf("NATS publisher: %T", natsPublisher)
	authRepo := repository.NewAuthRepository(db)
	logging.Logger.Debugf("Auth repo: %T", authRepo)
	sessionRepo := repository.NewSessionRepository(db)
	logging.Logger.Debugf("Session repo: %T", sessionRepo)
	roleRepo := repository.NewRoleRepository(db)
	logging.Logger.Debugf("Role repo: %T", roleRepo)
	logging.Logger.Debugf("Started repositories. AuthRepo: %T, SessionRepo: %T, RoleRepo: %T", authRepo, sessionRepo, roleRepo)

	logging.Logger.Debug("Starting services")
	// TODO: Add jwt settings to the config
	//jwtService := service.NewJwtService(defaultConfig.Jwt.Algo, defaultConfig.Jwt.Secret)
	jwtService := service.NewJwtService(jwt.SigningMethodHS256, "This is a secret key temp for avoid errors")

	roleService := service.NewRoleService(roleRepo)
	sessionService := service.NewSessionService(sessionRepo)
	authService := service.NewAuthService(authRepo, sessionService, jwtService, natsPublisher)
	logging.Logger.Debugf("Started services. Auth: %T, Session: %T, Role: %T", authService, sessionService, roleService)

	logging.Logger.Debug("Starting controllers")
	authAPI := api.NewAuthAPI(authService, sessionService, roleService)

	logging.Logger.Debug("Starting routes")
	router := r.Group("/")
	authAPI.RegisterRoutes(router)

	logging.Logger.Debug("Starting server")
	err = r.Run(cfg.Backend.ListenAddress + ":" + strconv.Itoa(cfg.Backend.ListenPort))

	if err != nil {
		logging.Logger.WithError(err).Error("Failed to start server")
		panic(err) // The server failed to start
	}
}
