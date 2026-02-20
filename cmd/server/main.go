package main

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/database"
	"github.com/padapook/bestbit-core/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.GormConnectDB(); err != nil {
		log.Fatal("[postgres gorm] Error connecting to database:", err)
	}

	if err := database.AutoMigrate(database.GormDB); err != nil {
		log.Fatal("[postgres gorm] Migration failed:", err)
	}

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:8081",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Routes(app, database.GormDB)

	app.Run(":" + os.Getenv("PORT"))
}
