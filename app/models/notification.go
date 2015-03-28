package models

import (
	"fmt"
	"time"
)

type Notification struct {
	Id             int32
	NotifyType     string `sql:"not null"`
	Read           bool   `sql:"default: false;not null"`
	UserId         int32  `sql:"not null"`
	User           User
	ActorId        int32 `sql:"not null"`
	Actor          User
	NotifyableType string `sql:"not null"`
	NotifyableId   int32  `sql:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NotifyInfo struct {
	UnreadCount int    `json:"unread_count"`
	Title       string `json:"title"`
	Avatar      string `json:"avatar"`
	Path        string `json:"path"`
}

func (n *Notification) Topic() (t Topic) {
	if n.NotifyableType == "Topic" {
		err := DB.Unscoped().First(&t, n.NotifyableId).Error
		if err != nil {
			return
		}
	}
	return
}

func (n *Notification) Reply() (r Reply) {
	if n.NotifyableType == "Reply" {
		err := DB.Unscoped().First(&r, n.NotifyableId).Error
		if err != nil {
			return
		}
	}
	return
}

func (n *Notification) NotifyableTitle() string {
	switch n.NotifyableType {
	case "Topic":
		return n.Topic().Title
	case "Reply":
		t := Topic{}
		db.First(&t, n.Reply().TopicId)
		return t.Title
	default:
		return ""
	}
}

func (n *Notification) NotifyableURL() string {
	switch n.NotifyableType {
	case "Topic":
		return fmt.Sprintf("/topics/%v", n.NotifyableId)
	case "Reply":
		return fmt.Sprintf("/topics/%v", n.Reply().TopicId)
	default:
		return ""
	}
}

func createNotification(notifyType string, userId int32, actorId int32, notifyableType string, notifyableId int32) error {
	note := Notification{
		NotifyType:     notifyType,
		UserId:         userId,
		ActorId:        actorId,
		NotifyableType: notifyableType,
		NotifyableId:   notifyableId,
	}

	exitCount := 0
	db.Model(Notification{}).Where(
		"user_id = ? and actor_id = ? and notifyable_type = ? and notifyable_id = ?",
		userId, actorId, notifyableType, notifyableId).Count(&exitCount)
	if exitCount > 0 {
		return nil
	}

	err := db.Save(&note).Error

	go PushNotifyInfoToUser(userId, note)

	return err
}

func (r *Reply) NotifyReply() error {
	if r.NewRecord() {
		return nil
	}

	if r.Topic.NewRecord() {
		return nil
	}

	if r.Topic.UserId == r.UserId {
		return nil
	}

	return createNotification("Reply", r.Topic.UserId, r.UserId, "Reply", r.Id)
}

func NotifyMention(userId, actorId int32, notifyableType string, notifyableId int32) error {
	return createNotification("Mention", userId, actorId, notifyableType, notifyableId)
}

func (u User) NotificationsPage(page, perPage int) (notes []Notification, pageInfo Pagination) {
	pageInfo = Pagination{}
	pageInfo.Query = db.Model(&Notification{}).Preload("Actor")
	pageInfo.Query = pageInfo.Query.Where("user_id = ?", u.Id).Order("id desc")

	pageInfo.Path = "/notifications"
	pageInfo.PerPage = perPage
	pageInfo.Paginate(page).Find(&notes)
	return
}

func (u User) ReadNotifications(notes []Notification) error {
	ids := []int32{}
	for _, note := range notes {
		ids = append(ids, note.Id)
	}
	if len(ids) > 0 {
		err := db.Model(Notification{}).Where("user_id = ? and id in (?)", u.Id, ids).Update("read", true).Error
		go PushNotifyInfoToUser(u.Id, Notification{})
		return err
	}

	return nil
}

func (u User) ClearNotifications() error {
	return db.Where("user_id = ?", u.Id).Delete(Notification{}).Error
}

func (n *Notification) IsTopic() bool {
	return n.NotifyType == "Topic"
}

func (n *Notification) IsReply() bool {
	return n.NotifyType == "Reply"
}

func (n *Notification) IsMention() bool {
	return n.NotifyType == "Mention"
}

func (n *Notification) IsNotifyableReply() bool {
	return n.NotifyableType == "Reply"
}

func (n *Notification) IsNotifyableTopic() bool {
	return n.NotifyableType == "Topic"
}
