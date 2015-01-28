package controllers

import (
	"fmt"
	//"fmt"
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Topics struct {
	App
}

func (c Topics) Index() revel.Result {
	return c.Result
}

func (c Topics) New() revel.Result {
	return c.Render("topics/new.html")
}

func (c Topics) Create() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}
	t := &Topic{Title: c.Params.Get("title"), Body: c.Params.Get("body")}

	t.UserId = currentUser.Id
	v := CreateTopic(t)
	if v.HasErrors() {
		return c.renderValidation("topics/new.html", v)
	}
	return c.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (c Topics) Show(id string) revel.Result {
	return c.RenderText(fmt.Sprintf("visit topic %v", id))
}
