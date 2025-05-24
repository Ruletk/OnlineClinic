package main

import (
	"doctor/internal/handler"
	"doctor/internal/repository"
	"doctor/internal/service"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1) init config
	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		panic(err) // The initialization of the app failed, without a config we can't do anything
	}
	logging.InitLogger(*cfg)

	logging.Logger.Debug("Setting up database connection")
	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to the database")
		panic(err) // Also, the initialization of the app failed, without a database we can't do anything
	}

	// 5) init layers
	repo := repository.NewDoctorRepository(db)
	svc := service.NewDoctorService(repo)
	h := handler.NewDoctorHandler(svc)

	// 6) gin router + routes
	r := gin.Default()
	h.RegisterRoutes(r)

	logging.Logger.Debug("Setting up NATS connection")
	logging.Logger.WithError(err).Error("Failed to connect to NATS. Disabling NATS features.")

}
