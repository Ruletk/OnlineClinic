package main

import (
	"auth/internal/api"
	nats2 "auth/internal/nats"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
	"errors"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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
		//panic(err) // Also, the initialization of the app failed, without a database we can't do anything
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	logging.Logger.Debug("Starting repositories")
	natsPublisher := nats2.NewPublisher(natsConn)
	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	redisStorage := repository.NewRedisStorage(redisClient, mainCtx)
	logging.Logger.Debugf("Started repositories. AuthRepo: %T, SessionRepo: %T, RoleRepo: %T", authRepo, sessionRepo, roleRepo)

	logging.Logger.Debug("Starting services")
	// TODO: Add jwt settings to the config
	//jwtService := service.NewJwtService(defaultConfig.Jwt.Algo, defaultConfig.Jwt.Secret)
	jwtService := service.NewJwtService(jwt.SigningMethodHS256, "This is a secret key temp for avoid errors")
	roleService := service.NewRoleService(roleRepo)
	sessionService := service.NewSessionService(sessionRepo)
	authService := service.NewAuthService(authRepo, sessionService, jwtService, natsPublisher, redisStorage)
	logging.Logger.Debugf("Started services. Auth: %T, Session: %T, Role: %T", authService, sessionService, roleService)

	logging.Logger.Debug("Starting controllers")
	authAPI := api.NewAuthAPI(authService, sessionService, roleService)
	logging.Logger.Debugf("Started controllers. AuthAPI: %T", authAPI)

	logging.Logger.Debug("Starting routes")
	router := r.Group("/")
	authAPI.RegisterRoutes(router)
	logging.Logger.Debug("Routes registered")

	logging.Logger.Debug("Starting server")

	srv := &http.Server{
		Addr:    cfg.Backend.ListenAddress + ":" + strconv.Itoa(cfg.Backend.ListenPort),
		Handler: r,
	}

	// Separate shutdown context
	shutdownCtx, shutdownCancel := context.WithTimeout(mainCtx, 10*time.Second)
	defer shutdownCancel()

	g, gCtx := errgroup.WithContext(shutdownCtx)
	g.Go(func() error {
		logging.Logger.Infof("Starting server on %s:%d", cfg.Backend.ListenAddress, cfg.Backend.ListenPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Logger.WithError(err).Error("Failed to start server")
			return err
		}
		return nil
	})

	// Shutdown goroutines

	g.Go(func() error {
		<-gCtx.Done()
		logging.Logger.Info("Shutting down gin server...")
		if err := srv.Shutdown(gCtx); err != nil {
			logging.Logger.WithError(err).Error("Failed to shutdown server")
			return err
		}
		logging.Logger.Info("Server shutdown complete")
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		logging.Logger.Info("Shutting down Redis connection...")
		if redisClient == nil {
			logging.Logger.Info("Redis client is not initialized, skipping shutdown")
			return nil
		}

		if err := redisClient.Close(); err != nil {
			logging.Logger.WithError(err).Error("Failed to close Redis connection")
			return err
		}
		logging.Logger.Info("Redis connection shutdown complete")
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		logging.Logger.Info("Shutting down NATS connection...")
		if natsConn == nil {
			logging.Logger.Info("NATS client is not initialized, skipping shutdown")
			return nil
		}

		if err := natsConn.Drain(); err != nil {
			logging.Logger.WithError(err).Error("Failed to flush NATS connection")
			return err
		}
		logging.Logger.Info("NATS connection shutdown complete")
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		logging.Logger.Info("Shutting down database connection...")
		db, err := db.DB()
		if err != nil {
			logging.Logger.WithError(err).Error("Failed to get database connection")
			return err
		}
		if err := db.Close(); err != nil {
			logging.Logger.WithError(err).Error("Failed to close database connection")
			return err
		}
		logging.Logger.Info("Database connection shutdown complete")
		return nil
	})

	if err := g.Wait(); err != nil {
		logging.Logger.WithError(err).Error("Error in main goroutine")
	}
}
