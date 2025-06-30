package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	UserKey = "user"
)

type (
	Middleware struct {
		router *gin.Engine
		db     *gorm.DB
	}
)

func NewMiddleware(router *gin.Engine, db *gorm.DB) *Middleware {
	return &Middleware{
		router: router,
		db:     db,
	}
}
