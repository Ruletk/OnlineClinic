package main

import (
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			Url:  fmt.Sprintf("http://%s:%s", "auth", "8080"),
		},
		"doctor": {
			Name: "doctor",
			Host: "doctor",
			Port: "8080",
			Url:  fmt.Sprintf("http://%s:%s", "doctor", "8080"),
		},
		"patient": {
			Name: "patient",
			Host: "patient",
			Port: "8080",
			Url:  fmt.Sprintf("http://%s:%s", "patient", "8080"),
		},
		"appointment": {
			Name: "appointment",
			Host: "appointment",
			Port: "8080",
			Url:  fmt.Sprintf("http://%s:%s", "appointment", "8080"),
		},
	}

	r.Use(middleware.PrometheusMiddleware())
	r.Use(middleware.TokenMiddleware())

	for serviceName, serviceConfig := range services {
		serviceGroup := r.Group(fmt.Sprintf("/%s", serviceName))
		{
			serviceGroup.Any("/*proxyPath", middleware.ReverseProxy(&serviceConfig))
		}
	}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if err := r.Run(":8080"); err != nil {
		logging.Logger.Fatalf("Failed to start server: %v", err)
		panic(err)
	}
}
