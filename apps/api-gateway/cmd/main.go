package main

import (
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	mainContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		panic(err) // The initialization of the app failed, without a config we can't do anything
	}
	logging.InitLogger(*cfg)

	_ = mainContext

	r := gin.Default()

	services := map[string]middleware.ServiceConfig{
		"auth": {
			Name: "auth",
			Host: "auth",
			Port: "8080",
		},
		"doctor": {
			Name: "doctor",
			Host: "doctor",
			Port: "8080",
		},
		"patient": {
			Name: "patient",
			Host: "patient",
			Port: "8080",
		},
		"appointment": {
			Name: "appointment",
			Host: "appointment",
			Port: "8080",
		},
	}
	r.Use(middleware.TokenMiddleware())

	for serviceName, serviceConfig := range services {
		serviceGroup := r.Group(fmt.Sprintf("/%s", serviceName))
		{
			// Для каждой группы применяем свой ProxyMiddleware
			serviceGroup.Any("/*proxyPath", middleware.ReverseProxy(&serviceConfig))
		}
	}

	if err := r.Run(":8080"); err != nil {
		logging.Logger.Fatalf("Failed to start server: %v", err)
		panic(err)
	}
}
