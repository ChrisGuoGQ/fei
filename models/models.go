package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func Bootstrap() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=chrisguo dbname=beego sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.AutoMigrate(&Group{}, &Camera{})
}
