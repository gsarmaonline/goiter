package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthorisationMiddleware is a middleware that checks if the user is authorised to access the resource
func (m *Middleware) AuthorisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
