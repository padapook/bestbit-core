package main

import (
    "github.com/padapook/bestbit-core/internal/database"
    "github.com/padapook/bestbit-core/internal/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    database.InitDB()

    r := gin.Default()
    routes.Routes(r)

    r.Run(":3000")
}