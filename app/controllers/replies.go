package controllers

import (
	"errors"
	"fmt"
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
)

type Replies struct {
	App
	topic Topic
}

func (c Replies) Create() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	reply := &Reply{Body: c.Params.Get("body")}
	err := DB.Where("id = ?", c.Params.Get("id")).First(&c.topic).Error
	if err != nil {
		return c.RenderError(err)
	}
	c.RenderArgs["topic"] = c.topic

	reply.TopicId = c.topic.Id
	reply.UserId = c.currentUser.Id
	v := CreateReply(reply)
	if v.HasErrors() {
		c.Flash.Error("回帖失败")
		return c.Redirect(fmt.Sprintf("/topics/%v", c.topic.Id))
	}
	return c.Redirect(fmt.Sprintf("/topics/%v#reply%v", c.topic.Id, c.topic.RepliesCount))
}

func (c Replies) Update() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	reply := &Reply{}
	err := DB.Model(reply).First(reply, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}
	if !c.isOwner(reply) {
		return c.RenderError(errors.New("Not allow."))
	}
	reply.Body = c.Params.Get("body")
	err = DB.Save(reply).Error
	if err != nil {
		return c.RenderError(err)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", reply.TopicId))
}

func (c Replies) Edit() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	reply := &Reply{}
	err := DB.Model(reply).First(reply, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}
	if !c.isOwner(reply) {
		return c.RenderError(errors.New("Not allow."))
	}
	c.RenderArgs["reply"] = reply
	return c.Render("replies/edit.html")
}

func (c Replies) Delete() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	reply := Reply{}
	err := DB.First(&reply, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}
	if !c.isOwner(reply) {
		return c.RenderError(errors.New("Not allow."))
	}

	DB.Delete(&reply)
	c.Flash.Success("回帖删除成功")
	return c.Redirect(fmt.Sprintf("/topics/%v", reply.TopicId))
}
