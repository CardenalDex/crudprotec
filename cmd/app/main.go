package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import your generated docs (you need to run `swag init` first)
	_ "github.com/CardenalDex/crudprotec/docs"

	config "github.com/CardenalDex/crudprotec/cmd"
	"github.com/CardenalDex/crudprotec/internal/adapter/handler"
	"github.com/CardenalDex/crudprotec/internal/adapter/repository"
	"github.com/CardenalDex/crudprotec/internal/usecase"
)

// @title Payment System API
// @version 1.0
// @description Clean Architecture Example with Go, Docker.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	db := repository.InitInternalDB()

	sqliteRepo := repository.NewSQLiteRepository(db)

	txService := usecase.NewTransactionService(sqliteRepo, sqliteRepo, sqliteRepo, sqliteRepo)

	txHandler := handler.NewTransactionHandler(txService)

	r := gin.Default()
	r.Use(config.RequestLoggerMiddleware())

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/transactions", txHandler.CreateTransaction)
	}

	log.Printf("Starting server on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
