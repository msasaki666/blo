package handlers

import (
	"app/models"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleListUser(c *gin.Context, db *gorm.DB) {
	var users []models.User
	if tx := db.Preload(clause.Associations).Find(&users); tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func HandleShowUser(c *gin.Context, db *gorm.DB) {
	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	var u models.User
	db.Preload(clause.Associations).First(&u, id)
	c.IndentedJSON(http.StatusOK, u)
}

func HandleCreateUser(c *gin.Context, db *gorm.DB) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	u := models.User{Email: email, PasswordDigest: string(passwordDigest)}
	db.Create(&u)
	tx := db.Preload(clause.Associations).First(&u)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, u)
}

// これまでのパスワードも確認する
func HandleUpdateUser(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	// 本人確認
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if user.ID != id {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	old := c.PostForm("old_password")
	// 前のパスワードをチェック
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(old)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("new_password")), bcryptCost)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	user.Email = c.PostForm("email")
	user.PasswordDigest = string(passwordDigest)
	tx := db.Clauses(clause.Returning{}).Save(&user)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	tx = db.Preload(clause.Associations).First(&user)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func HandleDeleteUser(c *gin.Context, db *gorm.DB, m *jwt.GinJWTMiddleware) {
	id, err := extractIDFromParam(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	// 本人確認
	user, err := currentUser(c, db, m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if user.ID != id {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	tx := db.Preload(clause.Associations).First(&user)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	// ポインタ渡しでいいのか
	tx = db.Delete(&user)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
