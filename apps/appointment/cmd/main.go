package main

import (
	controller2 "appointment/internal/controller"
	"appointment/internal/repository"
	"appointment/internal/service"
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		panic(err) // The initialization of the app failed, without a config we can't do anything
	}
	logging.InitLogger(*cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to Redis")
		panic(err) // The initialization of the app failed, without a Redis connection we can't do anything
	}

	logging.Logger.Info("Connected to Redis successfully")

	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to the database")
		panic(err) // The initialization of the app failed, without a database we can't do anything
	}

	appointmentRepo := repository.NewAppointmentRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)
	appointmentDBRedisRepo := repository.NewAppointmentDBRedisRepository(appointmentRepo, redisRepo, mainCtx)
	logging.Logger.Info("Appointment repository initialized successfully")

	appointmentService := service.NewAppointmentService(appointmentDBRedisRepo)

	controller := controller2.NewAppointmentController(appointmentService)

	r := gin.Default()
	r.Use(logging.GinLogger(logging.Logger), gin.Recovery())

	group := r.Group("")

	controller.RegisterRoutes(group)

	if err := r.Run(fmt.Sprintf(":8080")); err != nil {
		logging.Logger.WithError(err).Fatal("Failed to start the server")
	}
}
