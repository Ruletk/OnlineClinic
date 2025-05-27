package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres" // Или другой драйвер
	"gorm.io/gorm"

	grpcpackage "patient/internal/controllers/grpc"
	"patient/internal/controllers/rest"
	"patient/internal/models"
	proto "patient/internal/proto/gen"
	"patient/internal/repositories"
	"patient/internal/services"
)

func main() {
	// Инициализация БД (пример для PostgreSQL)
	dsn := "host=localhost user=postgres password=postgres dbname=clinic port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автомиграция для всех моделей
	db.AutoMigrate(
		&models.Patient{},
		&models.Allergy{},
		&models.Insurance{},
		&models.Prescription{},
	)

	// Инициализация слоев приложения
	patientRepo := repositories.NewPatientRepository(db)
	allergyRepo := repositories.NewAllergyRepository(db)
	insuranceRepo := repositories.NewInsuranceRepository(db)
	prescriptionRepo := repositories.NewPrescriptionRepository(db)

	allergyService := services.NewAllergyService(allergyRepo)
	insuranceService := services.NewInsuranceService(insuranceRepo)
	prescriptionService := services.NewPrescriptionService(prescriptionRepo)

	patientService := services.NewPatientService(
		patientRepo,
		allergyService,
		insuranceService,
		prescriptionService,
	)

	// Настройка REST сервера
	router := gin.Default()
	restHandler := rest.NewPatientHandler(patientService)
	restHandler.RegisterRoutes(router)
	restHandler.RegisterAllergyRoutes(router)
	restHandler.RegisterInsuranceRoutes(router)
	restHandler.RegisterPrescriptionRoutes(router)

	// Настройка gRPC сервера
	grpcServer := grpc.NewServer()
	grpcController := grpcpackage.NewPatientController(patientService)
	proto.RegisterPatientServiceServer(grpcServer, grpcController)

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск серверов
	go func() {
		log.Println("REST Server started on :8080")
		if err := router.Run(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server failed: %v", err)
		}
	}()

	go func() {
		lis, _ := net.Listen("tcp", ":50051")
		log.Println("gRPC Server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Остановка REST сервера
	if err := router.RunContext(ctx); err != nil {
		log.Println("REST shutdown error:", err)
	}

	// Остановка gRPC сервера
	grpcServer.GracefulStop()
	log.Println("Servers stopped gracefully")
}
