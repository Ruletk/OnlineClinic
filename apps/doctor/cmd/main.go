package main

import (
	"doctor/internal/handler"
	"doctor/internal/model"
	"doctor/internal/repository"
	"doctor/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1) load .env if exists
	_ = godotenv.Load()

	// 2) get DATABASE_URL, PORT
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3) connect to Postgres via GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 4) auto-migrate Doctor model
	if err := db.AutoMigrate(&model.Doctor{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// 5) init layers
	repo := repository.NewDoctorRepository(db)
	svc := service.NewDoctorService(repo)
	h := handler.NewDoctorHandler(svc)

	// 6) gin router + routes
	r := gin.Default()
	h.RegisterRoutes(r)

	// 7) start server
	log.Printf("starting doctor service on :%s …", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
