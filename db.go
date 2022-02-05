package main

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
