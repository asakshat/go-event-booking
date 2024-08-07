package main

import (
	"fmt"

	"time"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	initializers.InitKeys()
}

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080", "https://event-booking-app-real.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Server running  %s", "👨‍💻🏃"),
		})
	})
	routes.AuthRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	r.Run()
}
