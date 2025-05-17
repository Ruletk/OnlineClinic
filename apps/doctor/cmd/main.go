package main

import (
	"doctor/config"
	"doctor/internal/handler"
	"doctor/internal/repository"
	"doctor/internal/usecase"
	"doctor/pkg/database"
	"doctor/pkg/logger"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env
	config.Load()

	// Init logger
	logger.Init()
	log := logger.Log

	// Connect to DB
	pool := database.NewPostgresPool()
	defer pool.Close()

	// Init dependencies
	repo := repository.NewRepository(pool)
	uc := usecase.NewUseCase(repo)
	h := handler.NewHandler(uc)

	// Start server
	e := echo.New()
	h.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Starting doctor service on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal().Err(err).Msg("server stopped with error")
	}
}
