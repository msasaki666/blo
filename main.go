package main

import (
	"app/models"
	"app/routes"
	"log"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := setupRouter()
	r.Run()
}

func setupRouter() *gin.Engine {
	db := setupDb(models.MigrationTargets()...)
	// セキュアなトークンの生成方法
	// https://qiita.com/catatsuy/items/e21a889d52041e432d87
	secretKey, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		panic("Set SECRET_KEY environment variable")
	}

	tokenTimeout, ok := os.LookupEnv("TOKEN_TIMEOUT")
	if !ok {
		panic("Set SECRET_KEY environment variable")
	}

	d, err := time.ParseDuration(tokenTimeout)
	if err != nil {
		panic("time.ParseDuration()")
	}

	authMiddleware, err := createJwtMiddleware(secretKey, "email", db, d)
	if err != nil {
		panic("failed to create jwt middleware")
	}

	// When you use jwt.New(), the function is already automatically called for checking, which means you don't need to call it again.
	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{"/health"},
	}))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))
	// デフォルトで全てのプロキシを許可するからセキュリティ上よくないらしい
	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"status": "ok",
		})
	})

	routes.DrawUser(
		r.Group("/users"),
		db,
		authMiddleware,
	)

	routes.DrawPost(
		r.Group("/posts"),
		db,
		authMiddleware,
	)

	routes.DrawTag(
		r.Group("/tags"),
		db,
		authMiddleware,
	)
	return r
}
