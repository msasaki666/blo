package main

import (
	"app/app/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func main() {
	db := setupDb(&models.User{})
	// セキュアなトークンの生成方法
	// https://qiita.com/catatsuy/items/e21a889d52041e432d87
	secretKey, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		panic("Set SECRET_KEY environment variable")
	}
	authMiddleware, err := createJwtMiddleware(secretKey, "email", db)
	if err != nil {
		panic("failed to create jwt middleware")
	}
	// When you use jwt.New(), the function is already automatically called for checking, which means you don't need to call it again.
	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	ug := r.Group("/users")
	ug.POST("/login", authMiddleware.LoginHandler)

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
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
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
