package controllers

import (
	"fmt"
	"strconv"
	//"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Topics struct {
	App
}

func (c Topics) Index() revel.Result {
	offset, _ := strconv.Atoi(c.Params.Get("offset"))
	topics := FindTopicPages(offset, 20)
	c.RenderArgs["topics"] = topics
	return c.Render("topics/index.html")
}

func (c Topics) New() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{}
	c.RenderArgs["nodes"] = FindAllNodes()
	c.RenderArgs["topic"] = t
	return c.Render("topics/new.html")
}

func (c Topics) Create() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	nodeId, _ := strconv.Atoi(c.Params.Get("node_id"))
	t := &Topic{
		Title:  c.Params.Get("title"),
		Body:   c.Params.Get("body"),
		NodeId: int32(nodeId),
	}

	t.UserId = c.currentUser.Id
	v := CreateTopic(t)
	if v.HasErrors() {
		c.RenderArgs["topic"] = t
		c.RenderArgs["nodes"] = FindAllNodes()
		return c.renderValidation("topics/new.html", v)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (c Topics) Show() revel.Result {
	t := Topic{}
	DB.Preload("User").Preload("Node").First(&t, c.Params.Get("id"))
	replies := []Reply{}
	DB.Preload("User").Where("topic_id = ?", t.Id).Order("id asc").Find(&replies)
	c.RenderArgs["topic"] = t
	c.RenderArgs["replies"] = replies
	return c.Render("topics/show.html")
}

func (c Topics) Edit() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{}
	DB.Where("id = ?", c.Params.Get("id")).First(t)
	if !c.isOwner(t) {
		c.Flash.Error("没有修改的权限")
		return c.Redirect("/")
	}
	c.RenderArgs["topic"] = t
	c.RenderArgs["nodes"] = FindAllNodes()
	return c.Render("topics/edit.html")
}

func (c Topics) Update() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := Topic{}
	DB.First(&t, c.Params.Get("id"))
	if !c.isOwner(t) {
		c.Flash.Error("没有修改的权限")
		return c.Redirect("/")
	}
	nodeId, _ := strconv.Atoi(c.Params.Get("node_id"))
	t.NodeId = int32(nodeId)
	t.Title = c.Params.Get("title")
	t.Body = c.Params.Get("body")
	v := UpdateTopic(&t)
	if v.HasErrors() {
		c.RenderArgs["topic"] = t
		c.RenderArgs["nodes"] = FindAllNodes()
		return c.renderValidation("topics/edit.html", v)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (c Topics) Delete() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := Topic{}
	DB.First(&t, c.Params.Get("id"))
	if !c.isOwner(t) {
		c.Flash.Error("没有修改的权限")
		return c.Redirect("/")
	}

	err := DB.Delete(&t).Error
	if err != nil {
		c.RenderError(err)
	}
	return c.Redirect("/topics")
}
