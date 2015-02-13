package models

import (
	"github.com/revel/revel"
	"time"
)

type Topic struct {
	BaseModel
	UserId       int32 `sql:"not null"`
	User         User
	Title        string `sql:"size:300;not null"`
	Body         string `sql:"type:text;not null"`
	Replies      []Reply
	RepliesCount int32 `sql:"not null;default 0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Topic) validate() (v revel.Validation) {
	v = revel.Validation{}
	switch t.NewRecord() {
	case false:
	default:
		v.Required(t.UserId).Key("user_id").Message("不能为空")
		v.Min(int(t.UserId), 1).Key("user_id").Message("不正确")
		v.MinSize(t.Title, 10).Key("标题").Message("最少要 10 个子符")
		v.MaxSize(t.Title, 100).Key("标题").Message("最多只能写 100 个字符")
		v.MinSize(t.Body, 1).Key("内容").Message("不能为空")
	}
	return v
}

func FindTopicPages(offset, limit int) []Topic {
	topics := []Topic{}
	db.Preload("User").Order("id desc").Offset(offset).Limit(limit).Find(&topics)
	return topics
}

func CreateTopic(t *Topic) revel.Validation {
	v := t.validate()
	if v.HasErrors() {
		return v
	}

	err := db.Save(t).Error
	if err != nil {
		v.Error("服务器异常创建失败")
	}
	return v
}

func UpdateTopic(t *Topic) revel.Validation {
	v := t.validate()
	if v.HasErrors() {
		return v
	}

	err := db.Save(t).Error
	if err != nil {
		v.Error("服务器异常更新失败")
	}
	return v
}
