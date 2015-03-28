package controllers

import (
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
)

type Users struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Users).Before, revel.BEFORE)
//	revel.InterceptMethod((*Users).After, revel.AFTER)
//}

func (c Users) Show(login string) revel.Result {
	u, err := FindUserByLogin(login)
	if err != nil {
		return c.RenderError(err)
	}
	c.RenderArgs["user"] = u
	return c.Render("users/show.html")
}
