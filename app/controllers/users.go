package controllers

import (
	. "github.com/huacnlee/mediom/app/models"
	"github.com/revel/revel"
)

type Users struct {
	App
	user User
}

func init() {
	revel.InterceptMethod((*Users).Before, revel.BEFORE)
	// revel.InterceptMethod((*Users).After, revel.AFTER)
}

func (c *Users) Before() revel.Result {
	login := c.Params.Get("login")
	var err error
	c.user, err = FindUserByLogin(login)
	if err != nil {
		c.Finish(c.NotFound("页面不存在。"))
	}
	c.ViewArgs["user"] = c.user
	return nil
}

func (c Users) Show() revel.Result {
	recentTopics := []Topic{}
	DB.Order("id desc").Where("user_id = ?", c.user.Id).Limit(10).Find(&recentTopics)
	c.ViewArgs["recent_topics"] = recentTopics
	return c.Render()
}

func (c Users) Topics(login string) revel.Result {
	return c.Render()
}
