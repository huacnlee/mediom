package models

import (
	//"crypto/md5"
	//"encoding/hex"
	//"fmt"
	//"github.com/martini-contrib/binding"
	//"strings"
	"time"
)

type Topic struct {
	Id        int32
	UserId    int32  `sql:"not null"`
	Title     string `sql:"size:300;not null"`
	Body      string `sql:"type:text;not null"`
	Replies   []Reply
	CreatedAt time.Time
	UpdatedAt time.Time
}
