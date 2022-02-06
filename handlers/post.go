package handlers

import (
	"app/models"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleListPost(c *gin.Context, db *gorm.DB) {
	var posts []models.Post
	if tx := db.Preload(clause.Associations).Find(&posts); tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}
	c.IndentedJSON(http.StatusOK, posts)
}

func HandleShowPost(c *gin.Context, db *gorm.DB) {
	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	var p models.Post
	db.Preload(clause.Associations).First(&p, id)
	c.IndentedJSON(http.StatusOK, p)
}

func HandleCreatePost(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	post := models.Post{
		Title:   c.PostForm("title"),
		Content: c.PostForm("content"),
		UserID:  user.ID,
	}
	tx := db.Create(&post)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	tx = db.Preload(clause.Associations).First(&post)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func HandleUpdatePost(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
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

	var post models.Post
	tx := db.First(&post, id)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	if user.ID != post.UserID {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	post.Title = c.PostForm("title")
	post.Content = c.PostForm("content")
	tx = db.Clauses(clause.Returning{}).Save(&post)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func HandleDeletePost(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
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

	var post models.Post
	tx := db.First(&post, id)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	if user.ID != post.UserID {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	tx = db.Clauses(clause.Returning{}).Delete(&post)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusNotFound, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}
