package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1")
	{
		RegisterUserRoutes(v1, db)
	}
}
