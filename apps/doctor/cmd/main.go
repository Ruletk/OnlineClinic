package main

import (
    "fmt"

    "github.com/Ruletk/OnlineClinic/apps/doctor/internal/handler"
    "github.com/Ruletk/OnlineClinic/apps/doctor/internal/repository"
    "github.com/Ruletk/OnlineClinic/apps/doctor/internal/service"
    "github.com/Ruletk/OnlineClinic/pkg/config"
    "github.com/Ruletk/OnlineClinic/pkg/database"
    "github.com/Ruletk/OnlineClinic/pkg/logging"

    "github.com/gin-gonic/gin"
    "github.com/nats-io/nats.go"
)

func main() {
    // 1) Load configuration (from .env or real ENV)
    cfg, err := config.GetDefaultConfiguration()
    if err != nil {
        fmt.Printf("Error loading configuration: %v\n", err)
        panic(err)
    }

    // 2) Init structured logger
    logging.InitLogger(*cfg)

    // 3) Connect to Postgres via our pkg/database
    db, err := database.NewPostgresDatabase(cfg)
    if err != nil {
        logging.Logger.WithError(err).Fatal("Failed to connect to the database")
    }

    // 4) (Optional) NATS — if URL is set in cfg.Nats.URL
    var nc *nats.Conn
    if cfg.Nats.URL != "" {
        nc, err = nats.Connect(cfg.Nats.URL)
        if err != nil {
            logging.Logger.WithError(err).Error("Failed to connect to NATS. Continuing without NATS")
            nc = nil
        }
    }

    // 5) Wire up clean-arch layers
    repo := repository.NewDoctorRepository(db)
    svc := service.NewDoctorService(repo, nc)
    h := handler.NewDoctorHandler(svc)

    // 6) Setup Gin and register routes
    router := gin.Default()
    h.RegisterRoutes(router)

    // 7) Start HTTP server
    addr := fmt.Sprintf(":%s", cfg.Server.Port)
    logging.Logger.Infof("Starting doctor service on %s", addr)
    if err := router.Run(addr); err != nil {
        logging.Logger.WithError(err).Fatal("Doctor service stopped")
    }
}
