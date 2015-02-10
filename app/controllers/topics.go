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
	topics := []*Topic{}
	offset, _ := strconv.Atoi(c.Params.Get("offset"))
	DB.Order("id desc").Limit(20).Offset(offset).Find(&topics)
	c.RenderArgs["topics"] = topics
	return c.Render("topics/index.html")
}

func (c Topics) New() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{}
	c.RenderArgs["topic"] = t
	return c.Render("topics/new.html")
}

func (c Topics) Create() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{Title: c.Params.Get("title"), Body: c.Params.Get("body")}

	t.UserId = c.currentUser.Id
	v := CreateTopic(t)
	if v.HasErrors() {
		c.RenderArgs["topic"] = t
		return c.renderValidation("topics/new.html", v)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (c Topics) Show() revel.Result {
	t := &Topic{}
	DB.Where("id = ?", c.Params.Get("id")).First(t)
	c.RenderArgs["topic"] = t
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
	return c.Render("topics/edit.html")
}

func (c Topics) Update() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{}
	DB.Where("id = ?", c.Params.Get("id")).First(t)
	if !c.isOwner(t) {
		c.Flash.Error("没有修改的权限")
		return c.Redirect("/")
	}
	t.Title = c.Params.Get("title")
	t.Body = c.Params.Get("body")
	v := UpdateTopic(t)
	if v.HasErrors() {
		c.RenderArgs["topic"] = t
		return c.renderValidation("topics/edit.html", v)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}
