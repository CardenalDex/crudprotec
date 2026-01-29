package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/CardenalDex/crudprotec/docs"

	config "github.com/CardenalDex/crudprotec/cmd"
	"github.com/CardenalDex/crudprotec/internal/adapter/handler"
	"github.com/CardenalDex/crudprotec/internal/adapter/repository"
	"github.com/CardenalDex/crudprotec/internal/usecase"
)

// @title Payment System API
// @version 0.001
// @description Prueba tecnica Go,Docker,Gin,Swagger,CleanArquitecture, and so much more
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
	adService := usecase.NewAdminService(sqliteRepo, sqliteRepo)
	merchantService := usecase.NewMerchantService(sqliteRepo, sqliteRepo, sqliteRepo)

	txHandler := handler.NewTransactionHandler(txService)
	adminHandler := handler.NewAdminHandler(adService)

	merchantHandler := handler.NewMerchantHandler(merchantService)

	r := gin.Default()
	r.Use(config.RequestLoggerMiddleware())

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	v1trans := v1.Group("/transactions")
	{
		v1trans.POST("/new", txHandler.CreateTransaction)
		v1trans.GET("/:id", txHandler.GetTransaction)
		v1trans.GET("/bymerchant/:merchantID", txHandler.GetMerchantTransactions)
		v1trans.GET("/transactions", txHandler.GetAllTransactions)
	}

	admin := v1.Group("/admin")
	{
		admin.POST("/businesses/new", adminHandler.RegisterBusiness)
		admin.GET("/businesses/:id", adminHandler.GetBusiness)
		admin.PATCH("/businesses/:id/commission", adminHandler.UpdateBusinessCommission)
		admin.DELETE("/businesses/delete/:id", adminHandler.RemoveBusiness)
	}

	audit := v1.Group("/audit")
	{
		audit.GET("/", adminHandler.GetAllLogs)
		audit.GET("/:resource_id", adminHandler.GetAuditTrail)
	}
	merchants := v1.Group("/merchants")
	{
		merchants.POST("/new", merchantHandler.RegisterMerchant)
		merchants.GET("/:id", merchantHandler.GetMerchant)
		merchants.GET("/bybusiness/:businessID", merchantHandler.GetBusinessMerchants)
		merchants.DELETE("/delete/:id", merchantHandler.RemoveMerchant)
	}

	log.Printf("Starting server on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
