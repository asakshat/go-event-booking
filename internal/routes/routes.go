package routes

import (
	"github.com/asakshat/go-event-booking/internal/controllers"
	"github.com/asakshat/go-event-booking/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the routes for the application
func AuthRoutes(r *gin.Engine) {
	r.POST("/api/signup", controllers.SignUp)
	r.POST("/api/login", controllers.Login)
	r.POST("/api/verify-email", controllers.VerifyEmail)
	r.POST("/api/forgot-password", controllers.ForgetPassword)
	r.POST("/api/reset-password", controllers.ChangePassword)

	authRoutes := r.Group("/")
	authRoutes.Use(middlewares.Authenticate())
	{
		authRoutes.GET("/api/logged", controllers.GetLoggedUser)
		authRoutes.POST("/api/logout", controllers.Logout)
	}
}

func EventRoutes(r *gin.Engine) {
	r.GET("/api/event/:event_id", controllers.GetEventByID)
	r.GET("/api/event", controllers.GetAllEvents)

	eventRoutes := r.Group("/api/event")
	eventRoutes.Use(middlewares.Authenticate(), middlewares.RequireEmailVerified())
	{
		eventRoutes.POST("/create", controllers.CreateEvent)
		// eventRoutes.GET("/:id", controllers.GetEvent)
		eventRoutes.PUT("/update/:event_id", controllers.EditEvent)
		eventRoutes.DELETE("/delete/:event_id", controllers.DeleteEvent)
		eventRoutes.DELETE("/delete_perm/:event_id", controllers.DeletePerm)

		eventRoutes.PATCH("/undo-delete/:event_id", controllers.UndoDeleteEvent)
		eventRoutes.GET("/events_by_organizer", controllers.GetEventsByOrganizer)
	}
}

func TicketRoutes(r *gin.Engine) {
	r.POST("/api/ticket/:event_id", controllers.BuyTicket)
	r.POST("/api/ticket/verify/*token", controllers.VerifyTicket)

}
