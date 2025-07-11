package main

import (
	"log"
	"order-service/docs"
	"order-service/internal/app"

	"go.uber.org/zap"
)


var logger *zap.Logger


func init() {	
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Description = "Order Service API"
	docs.SwaggerInfo.Title = "Order Service API"
	docs.SwaggerInfo.Version = "1.0.0"
}


func main() {
	app, err := app.NewApp()
	if err != nil {
		logger.Fatal("Failed to create app", zap.Error(err))
		log.Fatal(err)
	}
	
	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start app", zap.Error(err))
	}
}