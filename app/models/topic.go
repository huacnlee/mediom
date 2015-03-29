package models

import (
	"errors"
	"fmt"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"time"
)

type Topic struct {
	BaseModel
	UserId             int32 `sql:"not null"`
	User               User
	NodeId             int32
	Node               Node
	Title              string `sql:"size:300;not null"`
	Body               string `sql:"type:text;not null"`
	Replies            []Reply
	RepliesCount       int32 `sql:"not null;default: 0"`
	LastActiveMark     int64 `sql:"not null; default: 0"`
	LastRepliedAt      time.Time
	LastReplyId        int32
	LastReplyUserId    int32
	LastReplyUser      User `sql:"size:255"`
	LastReplyUserLogin string
	StarsCount         int32 `sql:"not null; default: 0"`
	WatchesCount       int32 `sql:"not null; default: 0"`
	Rank               int32 `sql:"not null; default: 0"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

const (
	RankNoPoint = -1
	RankNormal  = 0
	RankAwesome = 1
)

func (t *Topic) BeforeCreate() (err error) {
	t.LastActiveMark = time.Now().Unix()
	return nil
}

func (t *Topic) AfterCreate() (err error) {
	go t.CheckMention()
	return nil
}

func (t *Topic) validate() (v revel.Validation) {
	v = revel.Validation{}
	switch t.NewRecord() {
	case false:
	default:
		v.Required(t.UserId).Key("user_id").Message("不能为空")
		v.Required(t.NodeId).Key("node_id").Message("不能为空")
		v.Min(int(t.UserId), 1).Key("user_id").Message("不正确")
		v.MinSize(t.Title, 10).Key("标题").Message("最少要 10 个子符")
		v.MaxSize(t.Title, 100).Key("标题").Message("最多只能写 100 个字符")
		v.MinSize(t.Body, 1).Key("内容").Message("不能为空")
	}
	return v
}

func FindTopicPages(channel string, nodeId, page, perPage int) (topics []Topic, pageInfo Pagination) {
	pageInfo = Pagination{}
	pageInfo.Query = db.Model(&Topic{}).Preload("User").Preload("Node")

	switch channel {
	case "recent":
		pageInfo.Query = pageInfo.Query.Order("id desc")
	case "popular":
		pageInfo.Query = pageInfo.Query.Where("rank = 1 or stars_count >= 5")
		pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	case "node":
		pageInfo.Query = pageInfo.Query.Where("node_id = ?", nodeId)
		pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	default:
		pageInfo.Query = pageInfo.Query.Where("rank >= 0").Order("last_active_mark desc, id desc")
	}
	pageInfo.Path = "/topics"
	pageInfo.PerPage = perPage
	pageInfo.Paginate(page).Find(&topics)
	return
}

func CreateTopic(t *Topic) revel.Validation {
	v := t.validate()
	if v.HasErrors() {
		return v
	}
	t.LastActiveMark = time.Now().Unix()
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

func (t *Topic) UpdateLastReply(reply *Reply) (err error) {
	if reply == nil {
		return errors.New("Reply is nil")
	}

	db.First(&reply.User, reply.UserId)
	err = db.Exec(`UPDATE topics SET updated_at = ?, last_active_mark = ?, last_replied_at = ?,
		last_reply_id = ?, last_reply_user_login = ?, last_reply_user_id = ? WHERE id = ?`,
		time.Now(),
		time.Now().Unix(),
		time.Now(),
		reply.Id,
		reply.User.Login,
		reply.UserId,
		reply.TopicId).Error

	return err
}

func (t Topic) UpdateRank(rank int) error {
	if t.NewRecord() {
		return errors.New("Give a empty record.")
	}

	return db.Model(t).Update("rank", rank).Error

}

func (t Topic) IsAwesome() bool {
	return t.Rank == RankAwesome
}

func (t Topic) IsNormal() bool {
	return t.Rank == RankNormal
}

func (t Topic) IsNoPoint() bool {
	return t.Rank == RankNoPoint
}

func (t Topic) URL() string {
	if t.NewRecord() {
		return ""
	}
	return fmt.Sprintf("%v/topics/%v", "https://127.0.0.1:3000", t.Id)
}

func (t Topic) FollowerIds() (ids []int32) {
	db.Model(Followable{}).Where("follow_type = 'Watch' and topic_id = ?", t.Id).Pluck("user_id", &ids)
	return
}

func TopicsCountCached() (count int) {
	if err := cache.Get("topics/total", &count); err != nil {
		if err = db.Model(Topic{}).Count(&count).Error; err == nil {
			go cache.Set("topics/total", count, 30*time.Minute)
		}
	}

	return
}
