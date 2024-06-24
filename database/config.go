package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func SetupDB() *gorm.DB {
	dsn := "thanakrit:neung7989@tcp(127.0.0.1:3306)/gin_rest_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	DB = db
	return DB
}
