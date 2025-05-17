package main

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"google.golang.org/grpc"
	"log"
	"net"
	"patient/internal/models"
	"patient/internal/repositories"
	"patient/internal/services"
)

func main() {
	// Инициализация конфига
	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		panic(err)
	}

	logging.InitLogger(*cfg)

	// Подключение к БД
	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		logging.Logger.Error("Failed to connect to database:", err)
		panic(err)
	}

	// Автомиграция
	if err := db.AutoMigrate(&models.Patient{}); err != nil {
		logging.Logger.Error("Failed to connect to database:", err)
		panic(err)
	}

	// Инициализация слоёв
	patientRepo := repositories.NewPatientRepository(db)
	patientService := services.NewPatientService(patientRepo)
	patientController := grpc.NewPatientController(patientService)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Logger.Error("Failed to connect to database:", err)
		panic(err)
	}

	s := grpc.NewServer()
	RegisterPatientServiceServer(s, patientController)

	log.Println("Patient Service running on :50051")
	if err := s.Serve(lis); err != nil {
		logging.Logger.Error("Failed to connect to database:", err)
		panic(err)
	}
}
