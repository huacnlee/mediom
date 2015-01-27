package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db gorm.DB
var DB *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "monster:123123@/foo?charset=utf8&parseTime=True")
	DB = &db
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{}, &Topic{})
	db.Model(&User{}).AddUniqueIndex("index_on_login", "login")
}
