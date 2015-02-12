package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"reflect"
)

var db gorm.DB
var DB *gorm.DB

type BaseModel struct {
	Id int32
}

func (m BaseModel) NewRecord() bool {
	return m.Id <= 0
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "monster:123123@/foo?charset=utf8&parseTime=True")
	DB = &db
	if err != nil {
		panic(err)
	}

	fmt.Println("------- typeof ", reflect.TypeOf(BaseModel{}))
	fmt.Println("------- typeof ", reflect.TypeOf(Topic{}))

	db.LogMode(true)
	db.AutoMigrate(&User{}, &Topic{})
	db.Model(&User{}).AddUniqueIndex("index_on_login", "login")
	db.Model(&Topic{}).AddIndex("index_on_user_id", "user_id")
}
