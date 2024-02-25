package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	db = GetDB()
	db.AutoMigrate(&User{}, &Category{}, &Article{}, &Comment{}, &Log{})
}

func GetDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Error creating sqlite file:", err.Error())
	}

	return db
}
