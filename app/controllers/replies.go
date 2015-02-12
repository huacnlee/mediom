package controllers

import (
	"github.com/revel/revel"
	. "mediom/app/models"
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
		return c.errorsJSON(1, v.Errors)
	}
	return c.successJSON(reply)
}

func (c Replies) Update() revel.Result {
	return c.RenderText("h")
}

func (c Replies) Edit() revel.Result {
	return c.RenderText("h")
}
