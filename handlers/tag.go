package handlers

import (
	"app/models"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleListTag(c *gin.Context, db *gorm.DB) {
	var tags []models.Tag
	if tx := db.Preload(clause.Associations).Find(&tags); tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}
	c.IndentedJSON(http.StatusOK, tags)
}

func HandleShowTag(c *gin.Context, db *gorm.DB) {
	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	var t models.Tag
	db.Preload(clause.Associations).First(&t, id)
	c.IndentedJSON(http.StatusOK, t)
}

func HandleCreateTag(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	tag := models.Tag{Name: c.PostForm("name"), UserID: user.ID}
	tx := db.Create(&tag)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	tx = db.Preload(clause.Associations).First(&tag)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, tag)
}

func HandleUpdateTag(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	var tag models.Tag
	tx := db.Preload(clause.Associations).First(&tag, id)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	if user.ID != tag.UserID {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	tag.Name = c.PostForm("name")
	tx = db.Clauses(clause.Returning{}).Save(&tag)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, tag)
}

func HandleDeleteTag(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	var tag models.Tag
	tx := db.First(&tag, id)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	if user.ID != tag.UserID {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	tx = db.Clauses(clause.Returning{}).Delete(&tag)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, tag)
}
