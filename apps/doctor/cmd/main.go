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
	// 1) load config (.env or real ENV)
	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		panic(fmt.Errorf("config load: %w", err))
	}

	// 2) init logger
	logging.InitLogger(*cfg)

	// 3) connect Postgres
	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("db connect failed")
	}

	// 4) wire up layers
	repo := repository.NewDoctorRepository(db)
	svc := service.NewDoctorService(repo) // one arg: repo
	h := handler.NewDoctorHandler(svc)

	// 5) setup Gin + routes
	router := gin.Default()
	h.RegisterRoutes(router)

	// 6) optional NATS
	if cfg.Nats.Url != "" {
		nc, err := nats.Connect(cfg.Nats.Url)
		if err != nil {
			logging.Logger.WithError(err).Error("nats connect failed, continuing without it")
		} else {
			defer nc.Close()
			// pass nc into your service or handler if you add publish logic
			_ = nc
		}
	}

	// 7) run HTTP server on configured address
	addr := fmt.Sprintf("%s:%d", cfg.Backend.ListenAddress, cfg.Backend.ListenPort)
	logging.Logger.Infof("starting doctor service on %s", addr)
	if err := router.Run(addr); err != nil {
		logging.Logger.WithError(err).Fatal("server stopped unexpectedly")
	}
}
