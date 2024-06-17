package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User interface not found"})
		return 0, false
	}

	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return 0, false
	}

	userIDUint := uint(userIDFloat)
	return userIDUint, true
}
