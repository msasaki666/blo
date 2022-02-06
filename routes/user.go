package routes

import (
	"app/handlers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DrawUser(r *gin.RouterGroup, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	r.POST("/login", m.LoginHandler)
	r.GET("/", func(c *gin.Context) {
		handlers.HandleListUser(c, db)
	})
	r.GET("/:id", func(c *gin.Context) {
		handlers.HandleShowUser(c, db)
	})
	r.POST("/", func(c *gin.Context) {
		handlers.HandleCreateUser(c, db)
	})
	r.PUT("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleUpdateUser(c, db, m)
	})
	r.DELETE("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleDeleteUser(c, db, m)
	})
}
