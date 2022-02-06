package routes

import (
	"app/handlers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DrawTag(r *gin.RouterGroup, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	r.GET("/", func(c *gin.Context) {
		handlers.HandleListTag(c, db)
	})
	r.GET("/:id", func(c *gin.Context) {
		handlers.HandleShowTag(c, db)
	})
	r.POST("/", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleCreateTag(c, db, m)
	})
	r.PUT("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleUpdateTag(c, db, m)
	})
	r.DELETE("/:id", m.MiddlewareFunc(), func(c *gin.Context) {
		handlers.HandleDeleteTag(c, db, m)
	})
}
