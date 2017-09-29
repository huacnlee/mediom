package controllers

import (
	"github.com/revel/revel"
)

type Notifications struct {
	App
}

//func init() {
//	revel.InterceptMethod((*Home).Before, revel.BEFORE)
//	revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c Notifications) Index() revel.Result {
	c.requireUser()
	var page int
	c.Params.Bind(&page, "page")
	notes, pageInfo := c.currentUser.NotificationsPage(page, 8)
	c.currentUser.ReadNotifications(notes)
	c.ViewArgs["notifications"] = notes
	c.ViewArgs["page_info"] = pageInfo
	return c.Render()
}

func (c Notifications) Clear() revel.Result {
	c.requireUser()

	c.currentUser.ClearNotifications()
	return c.Redirect("/notifications")
}
