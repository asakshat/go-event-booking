package controllers

import (
	"net/http"

	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/gin-gonic/gin"
)

func BuyTicket(c *gin.Context) {

	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
