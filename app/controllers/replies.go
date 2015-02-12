package controllers

import (
	"errors"
	"fmt"
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
