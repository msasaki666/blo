package main

import (
	"app/app/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := setupDb(&models.User{})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	ug := r.Group("/users")
	ug.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		var u models.User
		db.First(&u, id)
		c.JSON(http.StatusOK, u)
	})
	ug.POST("/", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		// costは4~31,これ以外の値を指定したときのdefaultは10
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		us := &models.User{Email: email, PasswordDigest: string(hashed)}
		db.Create(us)
		var u models.User
		db.Take(&u, us)
		c.JSON(http.StatusOK, u)
	})
	r.Run()
}

func setupDb(models ...interface{}) *gorm.DB {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		panic("Set DB_DSN environment variable")
	}

	db, err := gorm.Open(
		// このDSNはURL指定してもいい
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}

	if err = db.AutoMigrate(models...); err != nil {
		panic("failed to automigrate")
	}
	return db
}
