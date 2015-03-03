package models

import (
	"fmt"
	"time"
)

type Followable struct {
	Id         int32
	FollowType string `sql:"size:20; not null"`
	TopicId    int32
	Topic      Topic
	UserId     int32
	User       User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u User) isFollowed(ftype string, t Topic) bool {
	var existCount int
	DB.Model(&Followable{}).Where("follow_type = ? and topic_id = ? and user_id = ?", ftype, t.Id, u.Id).Count(&existCount)
	if existCount > 0 {
		return true
	} else {
		return false
	}
}

func (u User) follow(ftype string, t Topic) bool {
	if t.NewRecord() || u.NewRecord() {
		return false
	}

	if u.isFollowed(ftype, t) {
		return false
	}

	follow := Followable{FollowType: ftype, TopicId: t.Id, UserId: u.Id}
	if DB.Save(&follow).Error != nil {
		return false
	}
	t.updateFollowCounter(ftype)
	return true
}

func (u User) unFollow(ftype string, t Topic) bool {
	if t.NewRecord() || u.NewRecord() {
		return false
	}

	if !u.isFollowed(ftype, t) {
		return false
	}

	DB.Where("follow_type = ? and topic_id = ? and user_id = ?", ftype, t.Id, u.Id).Delete(&Followable{})
	t.updateFollowCounter(ftype)
	return true
}

func (t Topic) updateFollowCounter(ftype string) {
	var count int
	DB.Model(&Followable{}).Where("follow_type = ? and topic_id = ?", ftype, t.Id).Count(&count)

	counterCacheKey := "watches_count"
	if ftype == "Star" {
		counterCacheKey = "stars_count"
	}

	err := DB.Model(t).Update(counterCacheKey, count).Error
	if err != nil {
		fmt.Println("WARNING: updateFollowCounter execute failed: ", err)
	}

}

func (u User) IsWatched(t Topic) bool {
	return u.isFollowed("Watch", t)
}

func (u User) Watch(t Topic) bool {
	return u.follow("Watch", t)
}

func (u User) UnWatch(t Topic) bool {
	return u.unFollow("Watch", t)
}

func (u User) IsStared(t Topic) bool {
	return u.isFollowed("Star", t)
}

func (u User) Star(t Topic) bool {
	return u.follow("Star", t)
}

func (u User) UnStar(t Topic) bool {
	return u.unFollow("Star", t)
}
