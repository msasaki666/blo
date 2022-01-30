package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	setupDb()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Run()
}

func setupDb(models ...interface{}) {
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
}
