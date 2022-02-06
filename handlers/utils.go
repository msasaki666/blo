package handlers

import (
	"app/models"
	"fmt"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const bcryptCost = 10

func extractIDFromParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		return 0, fmt.Errorf("invalid ID: %v", idStr)
	}
	return uint(idInt), nil
}

func currentUser(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) (models.User, error) {
	u := m.IdentityHandler(c)
	us, ok := u.(*models.User)
	if !ok {
		return models.User{}, fmt.Errorf("failed type assertion")
	}

	var user models.User
	tx := db.Preload(clause.Associations).Where(us).First(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}
