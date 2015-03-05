package controllers

import (
	"github.com/revel/revel"
	"strconv"
)

type Notifications struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c Notifications) Index() revel.Result {
	if r := c.requireUser(); r != nil {
		return r
	}

	page, _ := strconv.Atoi(c.Params.Get("page"))
	notes, pageInfo := c.currentUser.NotificationsPage(page, 10)
	c.currentUser.ReadNotifications(notes)
	c.RenderArgs["notifications"] = notes
	c.RenderArgs["page_info"] = pageInfo
	return c.Render("notifications/index.html")
}
