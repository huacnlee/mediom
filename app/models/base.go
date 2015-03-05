package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var db gorm.DB
var DB *gorm.DB

type BaseModel struct {
	Id        int32
	DeletedAt *time.Time
}

func (m BaseModel) NewRecord() bool {
	return m.Id <= 0
}

func (m BaseModel) Destroy() error {
	err := db.Delete(&m).Error
	return err
}

func (m BaseModel) IsDeleted() bool {
	return m.DeletedAt != nil
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "monster:123123@/mediom?charset=utf8&parseTime=True")
	DB = &db
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{}, &Topic{}, &Reply{}, &Node{}, &Followable{}, &Notification{})
	db.Model(&User{}).AddUniqueIndex("index_on_login", "login")
	db.Model(&Topic{}).AddIndex("index_on_user_id", "user_id")
	db.Model(&Topic{}).AddIndex("index_on_last_active_mark_deleted_at", "last_active_mark", "deleted_at")
	db.Model(&Topic{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(&Topic{}).AddIndex("index_on_rank", "rank")
	db.Model(&User{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(&Reply{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(&Followable{}).AddUniqueIndex("index_on_followable", "followable_type", "followable_id")
	db.Model(&Notification{}).AddIndex("index_on_user_id", "user_id")
	db.Model(&Notification{}).AddIndex("index_on_notifyable", "notifyable_type", "notifyable_id")
}
