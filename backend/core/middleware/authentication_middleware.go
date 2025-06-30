package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
)

// AuthenticationMiddleware is a middleware that checks if the user is authenticated
func (m *Middleware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session")
		if err != nil {
			c.JSON(401, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		var user models.User
		if err := m.db.Where("google_id = ?", session).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user in context for use in handlers
		c.Set(UserKey, &user)
		c.Next()
	}
}
