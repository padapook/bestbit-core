package controller

import (
	"net/http"
	// "log"
	"github.com/gin-gonic/gin"

	"github.com/padapook/bestbit-core/internal/database"
)

type TestController struct{}

func NewTestController() *TestController {
    return &TestController{}
}

func (ctrl *TestController) PlukPing(c *gin.Context) {
	var result int
	if err := database.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "conn db",
		"message": "pong",
		"db_test": result,
	})
}
