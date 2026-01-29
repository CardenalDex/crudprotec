package main

import (
	"github.com/CardenalDex/crudprotec/internal/adapter/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Inside Docker, /data/app.db should be a mounted volume
	db, err := gorm.Open(sqlite.Open("/data/app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the GORM models
	db.AutoMigrate(&repository.BusinessModel{}, &repository.TransactionModel{})

	return db
}
