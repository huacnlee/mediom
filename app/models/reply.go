package models

import (
	"github.com/huacnlee/revel"
	"time"
)

type Reply struct {
	BaseModel
	UserId    int32 `sql:"not null"`
	User      User
	TopicId   int32 `sql:"not null"`
	Topic     Topic
	Body      string `sql:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Reply) BeforeCreate() (err error) {
	err = db.Exec("update topics set replies_count = (replies_count + 1) where id = ?", r.TopicId).Error
	return err
}

func (r *Reply) BeforeDelete() (err error) {
	err = db.Exec("update topics set replies_count = (replies_count - 1) where id = ?", r.TopicId).Error
	return err
}

func (r *Reply) AfterCreate() (err error) {
	db.Model(r).Related(&r.Topic)
	err = r.Topic.UpdateLastReply(r)
	r.NotifyReply()
	r.CheckMention()
	return nil
}

func (r *Reply) validate() (v revel.Validation) {
	v = revel.Validation{}
	switch r.NewRecord() {
	case false:
	default:
		v.Required(r.TopicId).Key("topic_id").Message("不能为空")
		v.Min(int(r.TopicId), 1).Key("topic_id").Message("不能为空")
		v.Required(r.UserId).Key("user_id").Message("不能为空")
		v.Min(int(r.UserId), 1).Key("user_id").Message("不能为空")
		v.MinSize(r.Body, 1).Key("内容").Message("不能为空")
		v.MaxSize(r.Body, 10000).Key("内容").Message("最多不允许超过 10000 个子符")
	}
	return v
}

func CreateReply(r *Reply) revel.Validation {
	v := r.validate()
	if v.HasErrors() {
		return v
	}

	err := db.Save(r).Error
	if err != nil {
		v.Error("服务器异常创建失败")
	}
	return v
}
