package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

func GormConnectDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok", 
		host, port, user, password, dbname)

	appEnv := os.Getenv("APP_ENV")
	loggerLevel := logger.Info
	if appEnv == "production" {
		loggerLevel = logger.Error
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  loggerLevel,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("[postgres gorm] failed to open DB: %w", err)
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("[postgres gorm] failed to get DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	GormDB = dbConn
	log.Println("[postgres gorm] connected to database successfully!!")
	return nil
}