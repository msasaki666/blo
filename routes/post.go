package routes

import (
	"app/handlers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DrawPost(r *gin.RouterGroup, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	r.GET("/", func(c *gin.Context) {
		handlers.HandleListPost(c, db)
	})
	r.GET("/:id", func(c *gin.Context) {
		handlers.HandleShowPost(c, db)
	})
	r.POST("/", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleCreatePost(c, db, m)
	})
	r.PUT("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleUpdatePost(c, db, m)
	})
	r.DELETE("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleDeletePost(c, db, m)
	})
}
