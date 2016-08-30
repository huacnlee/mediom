package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"log"
	"os"
	"time"
)

var db *gorm.DB
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

func InitDatabase() {
	adapter := revel.Config.StringDefault("gorm.adapter", "mysql")
	databaseURI := revel.Config.StringDefault("gorm.database_uri", "")
	var err error
	db, err = gorm.Open(adapter, databaseURI)
	DB = db
	if err != nil {
		panic(err)
	}

	db.LogMode(false)
	logger = Logger{log.New(os.Stdout, "  ", 0)}
	db.SetLogger(logger)
	db.AutoMigrate(&User{}, &Topic{}, &Reply{}, &Node{}, &NodeGroup{}, &Followable{}, &Notification{}, &Setting{})
	db.Model(NodeGroup{}).AddIndex("index_on_sort", "sort")
	db.Model(Node{}).AddIndex("index_on_group_and_sort", "node_group_id", "sort")
	db.Model(User{}).AddUniqueIndex("index_on_login", "login")
	db.Model(Topic{}).AddIndex("index_on_user_id", "user_id")
	db.Model(Topic{}).AddIndex("index_on_last_active_mark_deleted_at", "last_active_mark", "deleted_at")
	db.Model(Topic{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(Topic{}).AddIndex("index_on_rank", "rank")
	db.Model(User{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(Reply{}).AddIndex("index_on_deleted_at", "deleted_at")
	db.Model(Notification{}).AddIndex("index_on_user_id", "user_id")
	db.Model(Notification{}).AddIndex("index_on_notifyable", "notifyable_type", "notifyable_id")
	db.Model(Setting{}).AddUniqueIndex("index_on_key", "key")
	db.LogMode(true)

	initPubsub()
}
