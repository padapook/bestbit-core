package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/controller"
)

func TestRoutes(rg *gin.RouterGroup) {
	testController := controller.NewTestController()
	testping := rg.Group("")

	{
        testping.GET("/testping", testController.PlukPing)
    }

}