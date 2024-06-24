package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a placeholder for actual authentication logic
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.Header.Get("user_id")
		// For real scenarios, decode token or session to get the user ID
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userIDInt, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		c.Set("user_id", uint(userIDInt))

		c.Next()
	}
}
