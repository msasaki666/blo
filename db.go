package main

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func setupDb(dial gorm.Dialector, models ...interface{}) *gorm.DB {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	db, err := gorm.Open(
		// このDSNはURL指定してもいい
		dial,
		&gorm.Config{Logger: logger},
	)
	if err != nil {
		panic("failed to connect database")
	}

	// DB作り直したい時にアンコメント
	// sqlDB, _ := db.DB()
	// sqlDB.Exec("DROP DATABASE IF EXISTS app")
	// sqlDB.Exec("CREAT DATABASE app")

	if err = db.AutoMigrate(models...); err != nil {
		panic("failed to automigrate")
	}
	// デバッグ用にテーブル構造確認
	// tables, _ := db.Migrator().GetTables()
	// fmt.Println("----------------------------------------------------------------")
	// fmt.Println(tables)
	// for _, v := range models {
	// 	columnTypes, _ := db.Migrator().ColumnTypes(v)
	// 	var columNames []string
	// 	for _, v := range columnTypes {
	// 		columNames = append(columNames, v.Name())
	// 	}
	// 	fmt.Println(columNames)
	// }
	// fmt.Println("----------------------------------------------------------------")
	return db
}
