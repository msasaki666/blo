package main

import (
	"app/app/models"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

// 参考
// https://taisablog.com/archives/go-jwt
func createJwtMiddleware(secretKey, identityKey string, db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(secretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// ログイン時にTokenを発行するLoginHandlerではAuthenticatorとPayloadFuncが呼ばれる
		// PayloadFuncはペイロードに含めるクレームを設定
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				// これでこんな感じの内容のJWTがクライアントに返される
				// header: {"alg":"HS256","typ":"JWT"}
				// payload: {"email":"test2@example.com","exp":1644053776,"orig_iat":16440501}
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		// 未編集
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Email: claims[identityKey].(string),
			}
		},
		// ログイン時にTokenを発行するLoginHandlerではAuthenticatorとPayloadFuncが呼ばれる
		// Authenticatorはログイン認証の為の関数
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var ul models.UserLogin
			if err := c.ShouldBind(&ul); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			var u models.User
			tx := db.First(
				&u,
				models.User{Email: ul.Email},
			)
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(ul.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &u, nil
		},
		// 未編集
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.User); ok && v.Email == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		// TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
